package database

import (
	"strings"

	"github.com/mirage208/redis-go/internal/connection"
	"github.com/mirage208/redis-go/internal/kvcache"
	"github.com/mirage208/redis-go/internal/resp"
)

// DB is the interface for redis style storage engine
type DB interface {
	Exec(client *connection.Connection, cmdLine [][]byte) resp.Reply
	AfterClientClose(c *connection.Connection)
	Close()
}

type ConcurrentDB struct {
	cache *kvcache.KVCache
}

type ExecFunc func(db *ConcurrentDB, args [][]byte) resp.Reply

func NewConcurrentDB() *ConcurrentDB {
	return &ConcurrentDB{
		cache: kvcache.NewKVCache(),
	}
}

func (db *ConcurrentDB) Exec(client *connection.Connection, cmdLine [][]byte) resp.Reply {
	if len(cmdLine) == 0 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'exec' command")
	}
	cmdName := strings.ToLower(string(cmdLine[0]))
	switch cmdName {
	case "multi":
		return resp.MakeErrorReply("ERR 'multi' command not supported in concurrent DB")
	case "exec":
		return resp.MakeErrorReply("ERR 'exec' command not supported in concurrent DB")
	case "discard":
		return resp.MakeErrorReply("ERR 'discard' command not supported in concurrent DB")
	case "watch":
		return resp.MakeErrorReply("ERR 'watch' command not supported in concurrent DB")
	default:
		return db.executeCommand(cmdName, cmdLine[1:])
	}
}

func (db *ConcurrentDB) AfterClientClose(c *connection.Connection) {
	// No specific actions needed for concurrent DB
}

func (db *ConcurrentDB) Close() {
	// No specific actions needed for concurrent DB
}

func (db *ConcurrentDB) executeCommand(cmdName string, args [][]byte) resp.Reply {
	cmd, exists := cmdTable[cmdName]
	if !exists {
		return resp.MakeErrorReply("ERR unknown command '" + cmdName + "'")
	}
	return cmd.executer(db, args)
}
