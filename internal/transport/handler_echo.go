package transport

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/mirage208/redis-go/common/logger"
	"github.com/mirage208/redis-go/pkg/sync/atomic"
	"github.com/mirage208/redis-go/pkg/sync/wait"
)

// EchoHandler echos received line to client, using for test
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

// MakeEchoHandler creates EchoHandler
func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient is client for EchoHandler, using for test
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close connection
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	err := c.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Handle echos received line to client
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}
	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		// may occur: client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		//logger.Info("sleeping")
		//time.Sleep(10 * time.Second)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key any, val any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
