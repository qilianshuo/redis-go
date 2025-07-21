package resp

import (
	"bytes"
	"io"
	"testing"
)

func TestParseStream(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Reply
	}{
		{
			name:     "Simple string",
			input:    "+OK\r\n",
			expected: MakeStatusReply("OK"),
		},
		{
			name:     "Simple error",
			input:    "-ERR unknown command\r\n",
			expected: MakeErrorReply("ERR unknown command"),
		},
		{
			name:     "Integer",
			input:    ":12345\r\n",
			expected: MakeIntegerReply(12345),
		},
		{
			name:     "Bulk string",
			input:    "$5\r\nhello\r\n",
			expected: MakeBulkReply([]byte("hello")),
		},
		{
			name:     "Array",
			input:    "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: MakeMultiBulkReply([][]byte{[]byte("foo"), []byte("bar")}),
		},
		{
			name:     "Null bulk string",
			input:    "$-1\r\n",
			expected: MakeNullBulkReply(),
		},
		{
			name:     "Empty array",
			input:    "*0\r\n",
			expected: MakeEmptyMultiBulkReply(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			payloads := ParseStream(reader)
			payload := <-payloads
			if payload.Err != nil {
				t.Errorf("Unexpected error: %v", payload.Err)
			}
			if !bytes.Equal(payload.Data.ToBytes(), tt.expected.ToBytes()) {
				t.Errorf("Expected %v, got %v", tt.expected, payload.Data)
			}
		})
	}
}

func BenchmarkParseStream(b *testing.B) {
	input := "+OK\r\n-ERR unknown command\r\n:12345\r\n$5\r\nhello\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$-1\r\n*0\r\n"
	reader := bytes.NewReader([]byte(input))

	for i := 0; i < b.N; i++ {
		_, err := reader.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatal(err)
			return
		}
		payloads := ParseStream(reader)
		for payload := range payloads {
			if payload.Err != nil && payload.Err != io.EOF {
				b.Fatal(payload.Err)
			}
		}
	}
}
