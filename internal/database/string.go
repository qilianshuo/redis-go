package database

import (
	"time"

	"github.com/mirage208/redis-go/internal/resp"
)

const (
	upsertPolicy = iota // upsert means insert if not exists, or update if exists
	insertPolicy        // insert means insert only if not exists
	updatePolicy        // update means update if exists
)

const unlimitedTTL int64 = 0

func setExecuter(db *ConcurrentDB, args [][]byte) resp.Reply {
	if len(args) < 2 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'set' command")
	}
	key := string(args[0])
	value := args[1]
	policy := updatePolicy
	ttl := unlimitedTTL
	if len(args) > 2 {
		// todo: handle expiration
	}

	entity := &DataEntity{
		Data: value,
	}

	var ok bool
	switch policy {
	case upsertPolicy:
		ok = db.cache.PutEntity(key, entity)
	case insertPolicy:
		ok = db.cache.PutIfAbsent(key, entity)
	case updatePolicy:
		ok = db.cache.PutIfExists(key, entity)
	default:
		return resp.MakeErrorReply("ERR unknown policy for 'set' command")
	}
	if ok {
		if ttl != unlimitedTTL {
			expireTime := time.Now().Add(time.Duration(ttl) * time.Millisecond)
			db.cache.Expire(key, expireTime)
		} else {
		}
		return resp.MakeOkReply()
	}
	return resp.MakeOkReply()
}
func getExecuter(db *ConcurrentDB, args [][]byte) resp.Reply {
	if len(args) < 1 {
		return resp.MakeErrorReply("ERR wrong number of arguments for 'get' command")
	}
	key := string(args[0])
	value, exists := db.cache.GetEntity(key)
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
