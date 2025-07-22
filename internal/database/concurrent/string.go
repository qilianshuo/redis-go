package concurrent

import "github.com/mirage208/redis-go/internal/resp"

func setExecuter(db *ConcurrentDB, args [][]byte) resp.Reply {
	if len(args) < 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'set' command")
	}
	key := string(args[0])
	value := args[1]
	if len(args) > 2 {
		// todo: handle expiration
	}
	db.data.Put(key, value)
	return resp.MakeOkReply()
}
func getExecuter(db *ConcurrentDB, args [][]byte) resp.Reply {
	if len(args) < 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'get' command")
	}
	key := string(args[0])
	value, exists := db.data.Get(key)
	if !exists {
		return resp.MakeNullBulkReply()
	}
	return resp.MakeBulkReply(value.([]byte))
}
func init() {
	// Register all commands
	registerCommand("set", setExecuter)
	registerCommand("get", getExecuter)
}
