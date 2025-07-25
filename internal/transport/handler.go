package transport

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/mirage208/redis-go/common/logger"
	"github.com/mirage208/redis-go/internal/connection"
	"github.com/mirage208/redis-go/internal/database"
	"github.com/mirage208/redis-go/internal/database/concurrent"
	"github.com/mirage208/redis-go/internal/resp"
	"github.com/mirage208/redis-go/pkg/sync/atomic"
)

// RespHandler implements transport.Handler and serves as a redis service
type RespHandler struct {
	// TODO
	activeConn sync.Map // *client -> placeholder
	db         database.DB
	closing    atomic.Boolean // refusing new client and new request
}

func NewHandler() *RespHandler {
	// TODO
	// db := database.NewStandaloneServer()
	db := concurrent.NewConcurrentDB()
	return &RespHandler{
		db: db,
	}
}

func (h *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	// TODO
	if h.closing.Get() || ctx.Done() != nil {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}

	client := connection.NewConn(conn)
	h.activeConn.Store(client, struct{}{})

	ch := resp.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if errors.Is(payload.Err, io.EOF) || errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(client)
				logger.Info("connection closed: ", client.RemoteAddr())
				return
			}
			continue
		}

		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}

		r, ok := payload.Data.(*resp.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}

		//cmdLine := ""
		//for _, arg := range r.Args {
		//	cmdLine += string(arg) + " "
		//}
		//logger.Debug(cmdLine)

		result := h.db.Exec(client, r.Args)
		if result != nil {
			_, _ = client.Write(result.ToBytes())
		} else {
			_, _ = client.Write(resp.MakeErrorReply("unknown").ToBytes())
		}
	}
}

func (h *RespHandler) Close() error {
	// TODO
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key any, val any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}
