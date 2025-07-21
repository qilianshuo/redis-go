package resp

import (
	"bytes"
	"strconv"
)

// Reply is the interface that wraps the ToBytes method.
//
// ToBytes returns the reply in RESP bytes.
type Reply interface {
	ToBytes() []byte
}

const CRLF = "\r\n"

/* ---- Const Reply ----*/
var (
	pongBytes           = []byte("+PONG" + CRLF)
	okBytes             = []byte("+OK" + CRLF)
	nullBulkBytes       = []byte("$-1" + CRLF)
	emptyMultiBulkBytes = []byte("*0" + CRLF)
	noBytes             = []byte("" + CRLF)
)

type (
	PongReply           struct{}
	OkReply             struct{}
	NullBulkReply       struct{}
	EmptyMultiBulkReply struct{}
	NoReply             struct{}
)

var (
	pongReply           = new(PongReply)
	okReply             = new(OkReply)
	nullBulkReply       = new(NullBulkReply)
	emptyMultiBulkReply = new(EmptyMultiBulkReply)
	noReply             = new(NoReply)
)

func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func (r *NoReply) ToBytes() []byte {
	return noBytes
}

func MakePongReply() *PongReply {
	return pongReply
}

func MakeOkReply() *OkReply {
	return okReply
}

func MakeNullBulkReply() *NullBulkReply {
	return nullBulkReply
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return emptyMultiBulkReply
}

func MakeNoReply() *NoReply {
	return noReply
}

/* ---- Status Reply ---- */
type StatusReply struct {
	Status string
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

/* ---- Bulk Reply ---- */
type BulkReply struct {
	Arg []byte
}

func (r *BulkReply) ToBytes() []byte {
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

/* ---- Multi Bulk Reply ---- */
type MultiBulkReply struct {
	Args [][]byte
}

func (r *MultiBulkReply) ToBytes() []byte {
	var buf bytes.Buffer
	//Calculate the length of buffer
	argLen := len(r.Args)
	argLenStr := strconv.Itoa(argLen)
	bufLen := 1 + len(argLenStr) + 2
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
	_, _ = buf.WriteString("*")
	_, _ = buf.WriteString(argLenStr)
	_, _ = buf.WriteString(CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			_, _ = buf.WriteString("$-1")
			_, _ = buf.WriteString(CRLF)
		} else {
			_, _ = buf.WriteString("$")
			_, _ = buf.WriteString(strconv.Itoa(len(arg)))
			_, _ = buf.WriteString(CRLF)
			//Write bytes,avoid slice of byte to string(slicebytetostring)
			_, _ = buf.Write(arg)
			_, _ = buf.WriteString(CRLF)
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

/* ---- Integer Reply ---- */
type IntegerReply struct {
	Code int64
}

func (r *IntegerReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

func MakeIntegerReply(code int64) *IntegerReply {
	return &IntegerReply{Code: code}
}

/* ---- Error Reply ---- */
type ErrorReply struct {
	Msg string
}

func (r *ErrorReply) ToBytes() []byte {
	return []byte("-" + r.Msg + CRLF)
}

func MakeErrorReply(msg string) *ErrorReply {
	return &ErrorReply{Msg: msg}
}
