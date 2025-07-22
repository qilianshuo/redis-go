package concurrent

import "strings"

type Command struct {
	name     string   // Command name
	executer ExecFunc // Function to execute the command
}

var cmdTable = make(map[string]*Command)

// RegisterCommand registers a new command with the command table
func registerCommand(name string, executer ExecFunc) {
	cmdTable[strings.ToLower(name)] = &Command{
		name:     name,
		executer: executer,
	}
}
