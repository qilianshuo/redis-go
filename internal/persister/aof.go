package persister

import (
	"time"

	"github.com/mirage208/redis-go/internal/kvcache"
	"github.com/mirage208/redis-go/internal/resp"
	"github.com/mirage208/redis-go/pkg/logger"
)

const (
	FsyncAlways   = "always"
	FsyncEverysec = "everysec"
	FsyncNo       = "no"
)

const (
	aofChanSize = 1 << 10
)

// aofPayload represents a single AOF payload.
type payload struct {
	cmdLine [][]byte
}

// SaveCommand saves a command to the AOF log.
func (p *Persister) SaveCommand(cmdLine [][]byte) {
	if p.aofChan == nil {
		return
	}

	p.aofChan <- &payload{
		cmdLine: cmdLine,
	}
}

// LoadAof loads the AOF log into the cache.
func (p *Persister) LoadAof() *kvcache.Cache {
	// TODO: Implement AOF loading logic
	return nil
}

// Fsync flushes the AOF buffer to disk.
func (p *Persister) Fsync() {
	if err := p.aofFile.Sync(); err != nil {
		logger.Errorf("failed to fsync AOF file: %v", err)
	}
}

// listenAof listens for AOF payloads and writes them to the AOF file.
func (p *Persister) listenAof() {
	go func() {
		for pd := range p.aofChan {
			p.writeAof(pd)
		}
	}()
}

// writeAof writes an AOF payload to the AOF file.
func (p *Persister) writeAof(pd *payload) {
	data := resp.MakeMultiBulkReply(pd.cmdLine).ToBytes()
	_, err := p.aofFile.Write(data)
	if err != nil {
		logger.Warnf("failed to write AOF file: %v", err)
	}
	if p.aofFsync == FsyncAlways {
		p.Fsync()
	}
}

// fsyncEverySec periodically flushes the AOF buffer to disk.
func (p *Persister) fsyncEverySec() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.Fsync()
			case <-p.ctx.Done():
				return
			}
		}
	}()
}
