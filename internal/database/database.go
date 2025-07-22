package database

import (
	"github.com/mirage208/redis-go/internal/connection"
	"github.com/mirage208/redis-go/internal/resp"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// DB is the interface for redis style storage engine
type DB interface {
	Exec(client *connection.Connection, cmdLine [][]byte) resp.Reply
	AfterClientClose(c *connection.Connection)
	Close()
}

// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity struct {
	Data any
}
