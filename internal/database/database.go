package database

import (
	"strings"

	"github.com/mirage208/redis-go/common/utils"
	"github.com/mirage208/redis-go/internal/config"
	"github.com/mirage208/redis-go/internal/connection"
	"github.com/mirage208/redis-go/internal/kvcache"
	"github.com/mirage208/redis-go/internal/persister"
	"github.com/mirage208/redis-go/internal/resp"
)

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
type Command struct {
	name     string
	args     [][]byte
	callback chan resp.Reply
}

type SequentialDB struct {
	cache *kvcache.Cache

	cmdCh chan *Command

	persister *persister.Persister
}

func NewSequentialDB() *SequentialDB {
	d := &SequentialDB{
		cache: kvcache.NewKVCache(),
		cmdCh: make(chan *Command, 1024),
	}
	if config.Properties.AppendOnly {
		validAof := utils.FileExists(config.Properties.AppendFilename)
		persister, err := persister.NewPersister(config.Properties.AppendFilename, config.Properties.AppendFsync)
		if err != nil {
			// Handle error
		}
		d.persister = persister
		if validAof {
			d.cache = persister.LoadAof()
		}
	}

	go d.handleCommands()
	return d
}

func (db *SequentialDB) Exec(client *connection.Connection, cmdLine [][]byte) resp.Reply {
	cmd := &Command{
		name:     strings.ToLower(string(cmdLine[0])),
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
		switch cmd.name {
		case "multi":
			cmd.callback <- resp.MakeErrorReply("ERR 'multi' command not supported in concurrent DB")
		case "exec":
			cmd.callback <- resp.MakeErrorReply("ERR 'exec' command not supported in concurrent DB")
		case "discard":
			cmd.callback <- resp.MakeErrorReply("ERR 'discard' command not supported in concurrent DB")
		case "watch":
			cmd.callback <- resp.MakeErrorReply("ERR 'watch' command not supported in concurrent DB")
		default:
			cmd.callback <- db.executeCommand(cmd.name, cmd.args)
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
