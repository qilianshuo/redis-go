package resp

import "testing"

func TestPongReply_ToBytes(t *testing.T) {
	reply := MakePongReplay()
	expected := []byte("+PONG\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestOkReply_ToBytes(t *testing.T) {
	reply := MakeOkReplay()
	expected := []byte("+OK\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestNullBulkReply_ToBytes(t *testing.T) {
	reply := MakeNullBulkReply()
	expected := []byte("$-1\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestEmptyMultiBulkReply_ToBytes(t *testing.T) {
	reply := MakeEmptyMultiBulkReply()
	expected := []byte("*0\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestNoReply_ToBytes(t *testing.T) {
	reply := &NoReply{}
	expected := []byte("")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestUnknownErrorReply_ToBytes(t *testing.T) {
	reply := MakeUnknownErrorReply()
	expected := []byte("-ERR unknown\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestStatusReply_ToBytes(t *testing.T) {
	reply := MakeStatusReply("READY")
	expected := []byte("+READY\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBulkReply_ToBytes(t *testing.T) {
	reply := MakeBulkReply([]byte("Hello, Redis!"))
	expected := []byte("$13\r\nHello, Redis!\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestMultiBulkReply_ToBytes(t *testing.T) {
	reply := MakeMultiBulkReply([][]byte{[]byte("Hello"), []byte("World")})
	expected := []byte("*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestIntReply_ToBytes(t *testing.T) {
	reply := MakeIntReply(42)
	expected := []byte(":42\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestErrorReply_ToBytes(t *testing.T) {
	reply := MakeErrorReply("Invalid argument")
	expected := []byte("-Invalid argument\r\n")
	actual := reply.ToBytes()
	if string(actual) != string(expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
