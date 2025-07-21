package connection

import (
	"net"
	"sync"
	"time"

	"github.com/qilianshuo/redis-go/common/sync/wait"
	"github.com/qilianshuo/redis-go/pkg/logger"
)

type Connection struct {
	// TODO
	conn net.Conn

	// wait until finish sending data, used for graceful shutdown
	sendingData wait.Wait
}

var connPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

// NewConn creates Connection instance
func NewConn(conn net.Conn) *Connection {
	// TODO
	c, ok := connPool.Get().(*Connection)
	if !ok {
		logger.Error("connection pool make wrong type")
		return &Connection{
			conn: conn,
		}
	}
	c.conn = conn
	return c
}

// Write sends response to client over tcp client
func (c *Connection) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	c.sendingData.Add(1)
	defer func() {
		c.sendingData.Done()
	}()

	return c.conn.Write(b)
}

// Close disconnect with the client
func (c *Connection) Close() error {
	c.sendingData.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()

	connPool.Put(c)
	return nil
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) Name() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return ""
}
