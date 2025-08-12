package logger

import (
	"os"
	"strings"
	"testing"
)

func TestNewStdoutLogger(t *testing.T) {
	logger := NewStdoutLogger()
	if logger == nil {
		t.Fatal("NewStdoutLogger returned nil")
	}
	logger.Output(INFO, 1, "stdout logger test")
}

func TestNewFileLoggerAndSetup(t *testing.T) {
	dir := "testlogdir"
	_ = os.RemoveAll(dir)
	settings := &Settings{
		Path:       dir,
		Name:       "test",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	}
	logger, err := NewFileLogger(settings)
	if err != nil {
		t.Fatalf("NewFileLogger error: %v", err)
	}
	logger.Output(INFO, 1, "file logger test")
	Setup(settings)
	Info("setup info test")
	Infof("setup infof test: %d", 123)
	Warn("setup warn test")
	Warnf("setup warnf test: %s", "abc")
	Error("setup error test")
	Errorf("setup errorf test: %v", "err")
	Debug("setup debug test")
	Debugf("setup debugf test: %d", 456)
	Fatal("setup fatal test")
	Fatalf("setup fatalf test: %d", 789)
	// 检查日志文件是否生成
	files, _ := os.ReadDir(dir)
	found := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".log") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Log file not created in %s", dir)
	}
	_ = os.RemoveAll(dir)
}

func TestLoggerOutputFormat(t *testing.T) {
	logger := NewStdoutLogger()
	// 检查格式
	logger.Output(DEBUG, 1, "debug msg")
	logger.Output(INFO, 1, "info msg")
	logger.Output(WARNING, 1, "warn msg")
	logger.Output(ERROR, 1, "error msg")
	logger.Output(FATAL, 1, "fatal msg")
}
