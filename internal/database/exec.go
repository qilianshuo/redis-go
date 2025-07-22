package database

import (
	"strconv"
	"time"

	"github.com/qilianshuo/redis-go/internal/database/datastruct/list"
	"github.com/qilianshuo/redis-go/internal/resp"
)

type CommandProc func(*Server, [][]byte) resp.Reply

var cmdTable = map[string]CommandProc{
	"PING":   cmdPing,
	"GET":    cmdGet,
	"SET":    cmdSet,
	"TTL":    cmdTTL,
	"RPUSH":  cmdRPush,
	"LPUSH":  cmdLPush,
	"LRANGE": cmdLRange,
	"LLEN":   cmdLLen,
	"LINDEX": cmdLIndex,
	"LPOP":   cmdLPop,
	"RPOP":   cmdRPop,
	"SAVE":   cmdSave,
	"LOAD":   cmdLoad,
	// Add more commands as needed
}

// CommandPing handles the PING command
func cmdPing(db *Server, args [][]byte) resp.Reply {
	if len(args) == 0 {
		return resp.MakePongReply()
	} else if len(args) == 1 {
		return resp.MakeStatusReply(string(args[0]))
	} else {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'ping' command")
	}
}

// CommandGet handles the GET command
func cmdGet(db *Server, args [][]byte) resp.Reply {
	if len(args) != 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'get' command")
	}

	keyStr := string(args[0])

	// Check if key is expired
	if db.isExpired(keyStr) {
		// Delete the key if expired
		db.deleteExpiredKey(keyStr)
		return resp.MakeNullBulkReply()
	}

	// Get the value
	val, ok := db.data.Get(keyStr)
	if !ok {
		return resp.MakeNullBulkReply()
	}
	valString, _ := val.(string)
	valBytes := []byte(valString)
	return resp.MakeBulkReply(valBytes)

}

// CommandSet handles the SET command
func cmdSet(db *Server, args [][]byte) resp.Reply {
	if len(args) < 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'set' command")
	}

	keyStr := string(args[0])
	valStr := string(args[1])

	// Remove expiration if any
	db.ttlMap.Remove(keyStr)

	// Store in dictionary
	db.data.Put(keyStr, valStr)

	// Parse optional arguments (EX, PX, NX, XX)
	for i := 2; i < len(args); i++ {
		arg := string(args[i])
		if (arg == "EX" || arg == "ex") && i+1 < len(args) {
			// Set expiration in seconds
			seconds, _ := strconv.Atoi(string(args[i+1]))
			expiryTime := time.Now().Add(time.Duration(seconds) * time.Second)
			db.setExpiry(keyStr, expiryTime)
			i++ // Skip the next argument as we've used it
		} else if arg == "PX" || arg == "px" {
			seconds, _ := strconv.Atoi(string(args[i+1]))
			seconds = seconds / 1000
			expiryTime := time.Now().Add(time.Duration(seconds) * time.Second)
			db.setExpiry(keyStr, expiryTime)
			i++
		}
		// TODO Handle other optional arguments...
	}

	return resp.MakeOkReply()
}

// CommandTTL handles the TTL command
func cmdTTL(db *Server, args [][]byte) resp.Reply {
	if len(args) != 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'ttl' command")
	}

	keyStr := string(args[0])

	// Check if key exists
	if _, ok := db.data.Get(keyStr); !ok {
		return resp.MakeIntegerReply(-2)
	}

	// Get expiry
	if expireVal, ok := db.ttlMap.Get(keyStr); !ok {
		return resp.MakeIntegerReply(-1)
	} else {
		// Convert to time
		expiryTime := expireVal.(time.Time)
		// Calculate remaining time
		ttl := time.Until(expiryTime).Seconds()
		if ttl < 0 {
			// Key has expired
			db.deleteExpiredKey(keyStr)
			return resp.MakeIntegerReply(-2)
		}
		return resp.MakeIntegerReply(int64(ttl))
	}
}

// isExpired checks if a key is expired
func (server *Server) isExpired(key string) bool {
	// Get expiry time
	expireVal, ok := server.ttlMap.Get(key)
	if !ok {
		return false
	}
	if expireVal == nil {
		return false
	}

	// Convert to time
	expiryTime := *expireVal.(*time.Time)

	// Check if expired
	return time.Now().After(expiryTime)
}

// deleteExpiredKey deletes an expired key
func (server *Server) deleteExpiredKey(key string) {
	server.data.Remove(key)
	server.ttlMap.Remove(key)
}

// setExpiry sets an expiry time for a key
func (server *Server) setExpiry(key string, when time.Time) {
	// Store the expiry time
	t := new(time.Time)
	*t = when
	server.ttlMap.Put(key, t)
}

// CommandRPush handles the RPUSH command
func cmdRPush(db *Server, args [][]byte) resp.Reply {
	if len(args) < 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'rpush' command")
	}

	key := string(args[0])
	val, exists := db.data.Get(key)

	var quickList *list.QuickList
	if !exists {
		quickList = list.NewQuickList()
		db.data.Put(key, quickList)
	} else {
		var ok bool
		quickList, ok = val.(*list.QuickList)
		if !ok {
			return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	// Add all values to the list
	for i := 1; i < len(args); i++ {
		quickList.Add(string(args[i]))
	}

	return resp.MakeIntegerReply(int64(quickList.Len()))
}

// CommandLPush handles the LPUSH command
func cmdLPush(db *Server, args [][]byte) resp.Reply {
	if len(args) < 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'lpush' command")
	}

	key := string(args[0])
	val, exists := db.data.Get(key)

	var quickList *list.QuickList
	if !exists {
		quickList = list.NewQuickList()
		db.data.Put(key, quickList)
	} else {
		var ok bool
		quickList, ok = val.(*list.QuickList)
		if !ok {
			return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	// Add all values to the list in reverse order
	for i := len(args) - 1; i >= 1; i-- {
		quickList.Insert(0, string(args[i]))
	}

	return resp.MakeIntegerReply(int64(quickList.Len()))
}

// CommandLRange handles the LRANGE command
func cmdLRange(db *Server, args [][]byte) resp.Reply {
	if len(args) != 3 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'lrange' command")
	}

	key := string(args[0])
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return resp.MakeErrorReply("ERR value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return resp.MakeErrorReply("ERR value is not an integer or out of range")
	}

	val, exists := db.data.Get(key)
	if !exists {
		return resp.MakeEmptyMultiBulkReply()
	}

	quickList, ok := val.(*list.QuickList)
	if !ok {
		return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Handle negative indices
	length := quickList.Len()
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop + 1
	} else {
		stop = stop + 1
	}

	// Adjust indices
	if start < 0 {
		start = 0
	}
	if stop > length {
		stop = length
	}
	if start >= stop {
		return resp.MakeEmptyMultiBulkReply()
	}

	// Get the range
	elements := quickList.Range(start, stop)
	args = make([][]byte, len(elements))
	for i, elem := range elements {
		args[i] = []byte(elem.(string))
	}
	return resp.MakeMultiBulkReply(args)
}

// CommandLLen handles the LLEN command
func cmdLLen(db *Server, args [][]byte) resp.Reply {
	if len(args) != 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'llen' command")
	}

	key := string(args[0])
	val, exists := db.data.Get(key)
	if !exists {
		return resp.MakeIntegerReply(0)
	}

	quickList, ok := val.(*list.QuickList)
	if !ok {
		return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return resp.MakeIntegerReply(int64(quickList.Len()))
}

// CommandLIndex handles the LINDEX command
func cmdLIndex(db *Server, args [][]byte) resp.Reply {
	if len(args) != 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'lindex' command")
	}

	key := string(args[0])
	index, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return resp.MakeErrorReply("ERR value is not an integer or out of range")
	}

	val, exists := db.data.Get(key)
	if !exists {
		return resp.MakeNullBulkReply()
	}

	quickList, ok := val.(*list.QuickList)
	if !ok {
		return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Handle negative index
	length := quickList.Len()
	if index < 0 {
		index = length + index
	}
	if index < 0 || index >= length {
		return resp.MakeNullBulkReply()
	}

	elem := quickList.Get(index)
	return resp.MakeBulkReply([]byte(elem.(string)))
}

// CommandLPop handles the LPOP command
func cmdLPop(db *Server, args [][]byte) resp.Reply {
	if len(args) != 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'lpop' command")
	}

	key := string(args[0])
	val, exists := db.data.Get(key)
	if !exists {
		return resp.MakeNullBulkReply()
	}

	quickList, ok := val.(*list.QuickList)
	if !ok {
		return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if quickList.Len() == 0 {
		return resp.MakeNullBulkReply()
	}

	elem := quickList.Remove(0)
	if quickList.Len() == 0 {
		db.data.Remove(key)
	}
	return resp.MakeBulkReply([]byte(elem.(string)))
}

// CommandRPop handles the RPOP command
func cmdRPop(db *Server, args [][]byte) resp.Reply {
	if len(args) != 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'rpop' command")
	}

	key := string(args[0])
	val, exists := db.data.Get(key)
	if !exists {
		return resp.MakeNullBulkReply()
	}

	quickList, ok := val.(*list.QuickList)
	if !ok {
		return resp.MakeErrorReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if quickList.Len() == 0 {
		return resp.MakeNullBulkReply()
	}

	elem := quickList.RemoveLast()
	if quickList.Len() == 0 {
		db.data.Remove(key)
	}
	return resp.MakeBulkReply([]byte(elem.(string)))
}

// CommandSave handles the SAVE command
func cmdSave(db *Server, args [][]byte) resp.Reply {
	if len(args) != 0 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'save' command")
	}

	if err := db.SaveRDB(); err != nil {
		return resp.MakeErrorReply("ERR saving failed: " + err.Error())
	}

	return resp.MakeOkReply()
}

// CommandLoad handles the LOAD command
func cmdLoad(db *Server, args [][]byte) resp.Reply {
	if len(args) != 0 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'load' command")
	}

	// 先保存当前数据
	if err := db.SaveRDB(); err != nil {
		return resp.MakeErrorReply("ERR saving current data failed: " + err.Error())
	}

	// 加载 RDB 文件
	if err := db.loadRDB(); err != nil {
		return resp.MakeErrorReply("ERR loading failed: " + err.Error())
	}

	return resp.MakeOkReply()
}
