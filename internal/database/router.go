package database

import (
	"strings"

	"github.com/mirage208/redis-go/internal/resp"
)

type DBCommand struct {
	name     string   // Command name
	executer ExecFunc // Function to execute the command
}

type ExecFunc func(db *SequentialDB, args [][]byte) resp.Reply

var cmdTable = make(map[string]*DBCommand)

// RegisterCommand registers a new command with the command table
func registerCommand(name string, executer ExecFunc) {
	cmdTable[strings.ToLower(name)] = &DBCommand{
		name:     name,
		executer: executer,
	}
}
