package resp

import (
	"bytes"
	"strconv"
)

// Reply is the interface of redis serialization protocol message
type Reply interface {
	ToBytes() []byte
}

var (
	// CRLF is the line separator of redis serialization protocol
	CRLF = "\r\n"
)

// PongReply is +PONG
type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

// ToBytes marshal redis.Reply
func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReplay() *PongReply {
	return &PongReply{}
}

// OkReply is +OK
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

// ToBytes marshal redis.Reply
func (r *OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReplay() *OkReply {
	return &OkReply{}
}

// NullBulkReply is empty string
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

// ToBytes marshal redis.Reply
func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

// MakeNullBulkReply creates a new NullBulkReply
func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// EmptyMultiBulkReply is an empty list
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

// ToBytes marshal redis.Reply
func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// NoReply respond nothing, for commands like subscribe
type NoReply struct{}

var noBytes = []byte("")

// ToBytes marshal redis.Reply
func (r *NoReply) ToBytes() []byte {
	return noBytes
}

type UnknownErrorReply struct{}

var unknownErrorReplyBytes = []byte("-ERR unknown\r\n")

func (r *UnknownErrorReply) ToBytes() []byte {
	return unknownErrorReplyBytes
}
func MakeUnknownErrorReply() *UnknownErrorReply {
	return &UnknownErrorReply{}

}

/* ---- Status Reply ---- */

// StatusReply stores a simple status string
type StatusReply struct {
	Status string
}

// ToBytes marshal redis.Reply
func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

// MakeStatusReply creates StatusReply
func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

/* ---- Bulk Reply ---- */

// BulkReply stores a binary-safe string
type BulkReply struct {
	Arg []byte
}

// ToBytes marshal redis.Reply
func (r *BulkReply) ToBytes() []byte {
	if r.Arg == nil {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

// MakeBulkReply creates  BulkReply
func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

/* ---- Multi Bulk Reply ---- */

// MultiBulkReply stores a list of string
type MultiBulkReply struct {
	Args [][]byte
}

// ToBytes marshal redis.Reply
func (r *MultiBulkReply) ToBytes() []byte {
	var buf bytes.Buffer
	//Calculate the length of buffer
	argLen := len(r.Args)
	bufLen := 1 + len(strconv.Itoa(argLen)) + 2
	for _, arg := range r.Args {
		if arg == nil {
			bufLen += 3 + 2
		} else {
			bufLen += 1 + len(strconv.Itoa(len(arg))) + 2 + len(arg) + 2
		}
	}
	//Allocate memory
	buf.Grow(bufLen)
	//Write string step by step,avoid concat strings
	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(argLen))
	buf.WriteString(CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1")
			buf.WriteString(CRLF)
		} else {
			buf.WriteString("$")
			buf.WriteString(strconv.Itoa(len(arg)))
			buf.WriteString(CRLF)
			buf.Write(arg)
			buf.WriteString(CRLF)
		}
	}
	return buf.Bytes()
}

// MakeMultiBulkReply creates MultiBulkReply
func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

/* ---- Int Reply ---- */

// IntReply stores an int64 number
type IntReply struct {
	Code int64
}

// ToBytes marshal redis.Reply
func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

// MakeIntReply creates IntReply
func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

/* ---- Error Reply ---- */

// ErrorReply is an error and redis.Reply
type ErrorReply struct {
	Status string
}

// ToBytes marshal redis.Reply
func (r *ErrorReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

// MakeErrorReply
func (r *ErrorReply) Error() string {
	return r.Status
}

// MakeErrorReply creates ErrorReply
func MakeErrorReply(status string) *ErrorReply {
	return &ErrorReply{
		Status: status,
	}
}
