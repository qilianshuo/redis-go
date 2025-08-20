package persister

import (
	"context"
	"os"

	"github.com/mirage208/redis-go/pkg/logger"
)

// Persister is responsible for persisting data to disk.
type Persister struct {
	ctx    context.Context
	cancel context.CancelFunc

	aofFileName string
	aofFsync    string
	aofFile     *os.File
	aofChan     chan *payload
}

// NewPersister creates a new Persister.
func NewPersister(aofFileName string, aofFync string) (*Persister, error) {
	persister := &Persister{
		aofFileName: aofFileName,
		aofFsync:    aofFync,
		aofChan:     make(chan *payload, aofChanSize),
	}

	ctx, cancel := context.WithCancel(context.Background())
	persister.ctx, persister.cancel = ctx, cancel

	aofFile, err := os.OpenFile(persister.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	persister.aofFile = aofFile

	persister.listenAof()
	if persister.aofFsync == FsyncEverysec {
		persister.fsyncEverySec()
	}
	return persister, nil
}

// Close closes the AOF file and cancels the context.
func (p *Persister) Close() error {
	if p == nil {
		return nil
	}
	if p.aofFile != nil {
		err := p.aofFile.Close()
		if err != nil {
			logger.Warnf("failed to close AOF file: %v", err)
			return err
		}
	}
	p.cancel()
	return nil
}
