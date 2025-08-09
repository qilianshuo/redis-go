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

type CMD struct {
	cmd      string
	args     [][]byte
	callback chan resp.Reply
}

type SequentialDB struct {
	cache *kvcache.KVCache

	cmdCh chan *CMD
}

type ExecFunc func(db *SequentialDB, args [][]byte) resp.Reply

func NewSequentialDB() *SequentialDB {
	d := &SequentialDB{
		cache: kvcache.NewKVCache(),
		cmdCh: make(chan *CMD, 1024),
	}
	go d.handleCommands()
	return d
}

func (db *SequentialDB) Exec(client *connection.Connection, cmdLine [][]byte) resp.Reply {
	cmd := &CMD{
		cmd:      strings.ToLower(string(cmdLine[0])),
		args:     cmdLine[1:],
		callback: make(chan resp.Reply),
	}
	db.cmdCh <- cmd
	return <-cmd.callback
}

func (db *SequentialDB) AfterClientClose(c *connection.Connection) {
	// No specific actions needed for sequential DB
}

func (db *SequentialDB) Close() {
	// No specific actions needed for sequential DB
}

func (db *SequentialDB) handleCommands() {
	for cmd := range db.cmdCh {
		switch cmd.cmd {
		case "multi":
			cmd.callback <- resp.MakeErrorReply("ERR 'multi' command not supported in concurrent DB")
		case "exec":
			cmd.callback <- resp.MakeErrorReply("ERR 'exec' command not supported in concurrent DB")
		case "discard":
			cmd.callback <- resp.MakeErrorReply("ERR 'discard' command not supported in concurrent DB")
		case "watch":
			cmd.callback <- resp.MakeErrorReply("ERR 'watch' command not supported in concurrent DB")
		default:
			cmd.callback <- db.executeCommand(cmd.cmd, cmd.args)
		}
	}
}

func (db *SequentialDB) executeCommand(cmdName string, args [][]byte) resp.Reply {
	cmd, exists := cmdTable[cmdName]
	if !exists {
		return resp.MakeErrorReply("ERR unknown command '" + cmdName + "'")
	}
	return cmd.executer(db, args)
}
