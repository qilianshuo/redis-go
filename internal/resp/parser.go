package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"runtime/debug"
	"strconv"

	"github.com/mirage208/redis-go/pkg/logger"
)

// Payload represents the parsed RESP payload
type Payload struct {
	Data Reply
	Err  error
}

// ParseStream reads data from io.Reader and send payloads through channel
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse(reader, ch)
	return ch
}

func parse(rawReader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err, string(debug.Stack()))
		}
	}()

	reader := bufio.NewReader(rawReader)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &Payload{Err: err}
			close(ch)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// there are some empty lines within replication traffic, ignore this error
			//protocolError(ch, "empty line")
			continue
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+': // Simple strings
			content := string(line[1:])
			ch <- &Payload{
				Data: MakeStatusReply(content),
			}
		case '-': // Simple error
			ch <- &Payload{
				Data: MakeErrorReply(string(line[1:])),
			}
		case ':': // Integers
			value, parseErr := strconv.ParseInt(string(line[1:]), 10, 64)
			if parseErr != nil {
				err := errors.New("protocol error: " + "illegal number " + string(line[1:]))
				ch <- &Payload{Err: err}
				continue
			}
			ch <- &Payload{
				Data: MakeIntegerReply(value),
			}
		case '$': // Bulk strings
			err = parseBulkString(line, reader, ch)
			if err != nil {
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
		case '*': // Arrays
			err = parseArray(line, reader, ch)
			if err != nil {
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
		default: // Multi bulk strings, split by ' '
			args := bytes.Split(line, []byte{' '})
			ch <- &Payload{
				Data: MakeMultiBulkReply(args),
			}
		}
	}
}

func parseBulkString(header []byte, reader *bufio.Reader, ch chan<- *Payload) error {
	strLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil || strLen < -1 {
		protocolError(ch, "illegal bulk string header: "+string(header))
		return nil
	} else if strLen == -1 {
		ch <- &Payload{
			Data: MakeNullBulkReply(),
		}
		return nil
	}
	body := make([]byte, strLen+2)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return err
	}
	ch <- &Payload{
		Data: MakeBulkReply(body[:len(body)-2]),
	}
	return nil
}

func parseArray(header []byte, reader *bufio.Reader, ch chan<- *Payload) error {
	strNums, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil || strNums < 0 {
		protocolError(ch, "illegal array header "+string(header[1:]))
		return nil
	} else if strNums == 0 {
		ch <- &Payload{
			Data: MakeEmptyMultiBulkReply(),
		}
		return nil
	}

	lines := make([][]byte, 0, strNums)
	for i := int64(0); i < strNums; i++ {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		length := len(line)
		if length < 4 || line[length-2] != '\r' || line[0] != '$' {
			protocolError(ch, "illegal bulk string header "+string(line))
			break
		}
		strLen, err := strconv.ParseInt(string(line[1:length-2]), 10, 64)
		if err != nil || strLen < -1 {
			protocolError(ch, "illegal bulk string length "+string(line))
			break
		} else if strLen == -1 {
			lines = append(lines, []byte{})
		} else {
			body := make([]byte, strLen+2)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return err
			}
			lines = append(lines, body[:len(body)-2])
		}
	}
	ch <- &Payload{
		Data: MakeMultiBulkReply(lines),
	}
	return nil
}

func protocolError(ch chan<- *Payload, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &Payload{Err: err}
}
