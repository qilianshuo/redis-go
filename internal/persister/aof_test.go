package persister

import (
	"context"
	"os"
	"testing"
	"time"
)

type testPersister struct {
	Persister
}

func newTestPersister() *testPersister {
	p := &testPersister{}
	p.aofChan = make(chan *payload, 10)
	tmpFile, _ := os.CreateTemp("", "aof_test_*.aof")
	p.aofFile = tmpFile
	p.aofFsync = FsyncAlways
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return p
}

func TestSaveCommand(t *testing.T) {
	p := newTestPersister()
	go p.listenAof()
	cmd := [][]byte{[]byte("SET"), []byte("key"), []byte("val")}
	p.SaveCommand(cmd)
	time.Sleep(10 * time.Millisecond)
	p.cancel()
	close(p.aofChan)
	fi, _ := p.aofFile.Stat()
	if fi.Size() == 0 {
		t.Errorf("AOF not written")
	}
	p.aofFile.Close()
	os.Remove(p.aofFile.Name())
}

func TestFsync(t *testing.T) {
	p := newTestPersister()
	p.Fsync()
	p.aofFile.Close()
	os.Remove(p.aofFile.Name())
}

func TestWriteAof(t *testing.T) {
	p := newTestPersister()
	pd := &payload{cmdLine: [][]byte{[]byte("PING")}}
	p.writeAof(pd)
	fi, _ := p.aofFile.Stat()
	if fi.Size() == 0 {
		t.Errorf("writeAof did not write")
	}
	p.aofFile.Close()
	os.Remove(p.aofFile.Name())
}

func TestFsyncEverySec(t *testing.T) {
	p := newTestPersister()
	p.fsyncEverySec()
	time.Sleep(1100 * time.Millisecond)
	p.cancel()
	p.aofFile.Close()
	os.Remove(p.aofFile.Name())
}
