package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/qilianshuo/redis-go/internal/connection"
	"github.com/qilianshuo/redis-go/internal/database/datastruct/dict"
	"github.com/qilianshuo/redis-go/internal/database/rdb"
	"github.com/qilianshuo/redis-go/internal/resp"
	"github.com/qilianshuo/redis-go/pkg/logger"
)

// Command represents a command to be executed by the database
type Command struct {
	Name   string          // Command name
	Args   [][]byte        // Command arguments
	RespCh chan resp.Reply // Response channel
}

type Server struct {
	data   *dict.SequentialDict
	ttlMap *dict.SequentialDict

	cmdCh chan *Command

	closeChan chan struct{}

	// RDB 相关配置
	rdbFilename  string
	saveInterval time.Duration
	lastSaveTime time.Time
}

func (server *Server) processCommands() {
	logger.Debug("Begin processing commands")
	for {
		select {
		case cmd := <-server.cmdCh:
			proc, exists := cmdTable[cmd.Name]
			var result resp.Reply
			if exists {
				result = proc(server, cmd.Args)
			} else {
				result = resp.MakeErrorReply("ERR unknown command '" + cmd.Name + "'")
			}

			// Send response
			if cmd.RespCh != nil {
				cmd.RespCh <- result
			}
		case <-server.closeChan:
			return
		}
	}
}

func NewStandaloneServer() *Server {
	server := &Server{
		data:   dict.NewSequentialDict(),
		ttlMap: dict.NewSequentialDict(),

		cmdCh: make(chan *Command),

		// 设置默认的 RDB 配置
		rdbFilename:  "dump.rdb",
		saveInterval: 5 * time.Minute,
		lastSaveTime: time.Now(),
	}

	// 尝试加载 RDB 文件
	if err := server.loadRDB(); err != nil {
		logger.Error("Failed to load RDB file:", err)
	}

	// 启动自动保存
	go server.autoSave()

	logger.Debug("Start db server")
	go server.processCommands()
	return server
}

// SaveRDB 保存数据库快照
func (server *Server) SaveRDB() error {
	// 准备要保存的数据
	saveData := struct {
		Data   map[string]interface{} `json:"data"`
		TTLMap map[string]time.Time   `json:"ttl_map"`
	}{
		Data:   make(map[string]interface{}),
		TTLMap: make(map[string]time.Time),
	}

	// 收集数据
	server.data.ForEach(func(key string, value interface{}) bool {
		saveData.Data[key] = value
		return true
	})

	server.ttlMap.ForEach(func(key string, value interface{}) bool {
		if t, ok := value.(time.Time); ok {
			saveData.TTLMap[key] = t
		}
		return true
	})

	// 保存到文件
	err := rdb.Save(saveData, server.rdbFilename)
	if err == nil {
		server.lastSaveTime = time.Now()
		logger.Info("RDB saved successfully")
	}
	return err
}

// loadRDB 加载数据库快照
func (server *Server) loadRDB() error {
	data, err := rdb.Load(server.rdbFilename)
	if err != nil {
		return err
	}

	// 解析数据
	saveData, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid RDB data format")
	}

	// 恢复数据
	if dataMap, ok := saveData["data"].(map[string]interface{}); ok {
		for k, v := range dataMap {
			server.data.Put(k, v)
		}
	}

	if ttlMap, ok := saveData["ttl_map"].(map[string]interface{}); ok {
		for k, v := range ttlMap {
			if t, ok := v.(time.Time); ok {
				server.ttlMap.Put(k, t)
			}
		}
	}

	logger.Info("RDB loaded successfully")
	return nil
}

// autoSave 自动保存数据库快照
func (server *Server) autoSave() {
	ticker := time.NewTicker(server.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := server.SaveRDB(); err != nil {
				logger.Error("Failed to auto save RDB:", err)
			}
		case <-server.closeChan:
			// 在关闭前保存一次
			if err := server.SaveRDB(); err != nil {
				logger.Error("Failed to save RDB before shutdown:", err)
			}
			return
		}
	}
}

func (server *Server) Exec(c *connection.Connection, cmdLine [][]byte) resp.Reply {
	if c == nil {
		return resp.MakeErrorReply("ERR operation not permitted")
	}
	cmd := &Command{
		Name:   strings.ToUpper(string(cmdLine[0])),
		Args:   cmdLine[1:],
		RespCh: make(chan resp.Reply),
	}
	server.cmdCh <- cmd
	return <-cmd.RespCh
}

func (server *Server) AfterClientClose(c *connection.Connection) {
	// TODO
	logger.Debug(fmt.Sprintf("client closed: %s", c.RemoteAddr()))
}

// Close graceful shutdown database
func (server *Server) Close() {
	// 触发关闭信号
	server.closeChan <- struct{}{}
}

// GetEntity returns DataEntity bind to given key
func (server *Server) GetEntity(key string) (*DataEntity, bool) {

	raw, ok := server.data.Get(key)
	if !ok {
		return nil, false
	}
	if server.IsExpired(key) {
		server.data.Remove(key)
		server.ttlMap.Remove(key)
		return nil, false
	}
	entity, _ := raw.(*DataEntity)
	return entity, true
}

// PutEntity a DataEntity into DB
func (server *Server) PutEntity(key string, entity *DataEntity) int {
	ret := server.data.Put(key, entity)
	return ret
}

// PutIfExists edit an existing DataEntity
func (server *Server) PutIfExists(key string, entity *DataEntity) int {
	return server.data.PutIfExists(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (server *Server) PutIfAbsent(key string, entity *DataEntity) int {
	ret := server.data.PutIfAbsent(key, entity)
	return ret
}

// Remove the given key from db
func (server *Server) Remove(key string) {
	server.data.Remove(key)
	server.ttlMap.Remove(key)

}

// Removes the given keys from db
func (server *Server) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := server.data.Get(key)
		if exists {
			server.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
// deprecated
// for test only
func (server *Server) Flush() {
	server.data.Clear()
	server.ttlMap.Clear()
}

/* ---- TTL Functions ---- */

// IsExpired check whether a key is expired
func (server *Server) IsExpired(key string) bool {
	rawExpireTime, ok := server.ttlMap.Get(key)
	if !ok {
		return false
	}
	expireTime, _ := rawExpireTime.(time.Time)
	expired := time.Now().After(expireTime)
	if expired {
		server.Remove(key)
	}
	return expired
}

// ForEach traverses all the keys in the database
func (server *Server) ForEach(cb func(key string, data *DataEntity, expiration *time.Time) bool) {
	server.data.ForEach(func(key string, raw interface{}) bool {
		entity, _ := raw.(*DataEntity)
		var expiration *time.Time
		rawExpireTime, ok := server.ttlMap.Get(key)
		if ok {
			expireTime, _ := rawExpireTime.(time.Time)
			expiration = &expireTime
		}

		return cb(key, entity, expiration)
	})
}
