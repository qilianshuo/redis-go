package resp

import (
	"reflect"
	"strconv"
	"testing"
)

func TestPongReply_ToBytes(t *testing.T) {
	r := &PongReply{}
	result := r.ToBytes()
	expected := []byte("+PONG\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, go t %v", expected, result)
	}
}

func TestOkReply_ToBytes(t *testing.T) {
	r := &OkReply{}
	result := r.ToBytes()
	expected := []byte("+OK\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestNullReply_ToBytes(t *testing.T) {
	r := &NullBulkReply{}
	result := r.ToBytes()
	expected := []byte("$-1\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestEmptyReply_ToBytes(t *testing.T) {
	r := &EmptyMultiBulkReply{}
	result := r.ToBytes()
	expected := []byte("*0\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestNoReply_ToBytes(t *testing.T) {
	r := &NoReply{}
	result := r.ToBytes()
	expected := []byte("\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStatusReply_ToBytes(t *testing.T) {
	status := "TestStatus"
	r := &StatusReply{Status: status}
	result := r.ToBytes()
	expected := []byte("+" + status + "\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestBulkReply_ToBytes(t *testing.T) {
	arg := []byte("test arg")
	r := &BulkReply{Arg: arg}
	result := r.ToBytes()
	expected := []byte("$" + strconv.Itoa(len(arg)) + "\r\n" + string(arg) + "\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMultiBulkReply_ToBytes(t *testing.T) {
	args := [][]byte{[]byte("arg1"), []byte("arg2")}
	r := &MultiBulkReply{Args: args}
	result := r.ToBytes()
	expected := []byte("*2\r\n$4\r\narg1\r\n$4\r\narg2\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试包含nil元素
	args = [][]byte{[]byte("foo"), nil, []byte("bar")}
	r = &MultiBulkReply{Args: args}
	result = r.ToBytes()
	expected = []byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试空数组
	r = &MultiBulkReply{Args: [][]byte{}}
	result = r.ToBytes()
	expected = []byte("*0\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestIntegerReply_ToBytes(t *testing.T) {
	code := int64(42)
	r := &IntegerReply{Code: code}
	result := r.ToBytes()
	expected := []byte(":" + strconv.FormatInt(code, 10) + "\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestErrorReply_ToBytes(t *testing.T) {
	msg := "TestError"
	r := &ErrorReply{Msg: msg}
	result := r.ToBytes()
	expected := []byte("-" + msg + "\r\n")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expecte d %v, got %v", expected, result)
	}
}
