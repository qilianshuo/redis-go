package persister

import (
	"os"

	"github.com/mirage208/redis-go/internal/config"
	"github.com/mirage208/redis-go/pkg/logger"
)

func (p *Persister) GenerateRDB(rdbFilename string) error {
	file, err := os.CreateTemp(config.GetTmpDir(), "rdb-*.rdb")
	if err != nil {
		logger.Warnf("failed to create RDB file: %v", err)
		return err
	}
	defer file.Close()

	// todo: implement RDB generation logic

	return nil
}
