package persister

import (
	"time"

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

type payload struct {
	cmdLine [][]byte
}

func (p *Persister) Fsync() {
	if err := p.aofFile.Sync(); err != nil {
		logger.Errorf("failed to fsync AOF file: %v", err)
	}
}

func (p *Persister) listenAof() {
	go func() {
		for pd := range p.aofChan {
			p.writeAof(pd)
		}
	}()
}

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
