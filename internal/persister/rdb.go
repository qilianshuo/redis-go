package persister

import (
	"os"
	"strconv"
	"time"

	"github.com/mirage208/redis-go/internal/config"
	"github.com/mirage208/redis-go/internal/kvcache"
	"github.com/mirage208/redis-go/pkg/logger"

	rdbEncoder "github.com/hdt3213/rdb/encoder"
)

// GenerateRDB generates an RDB file from the current dataset.
func (p *Persister) GenerateRDB(rdbFilename string) error {
	file, err := os.CreateTemp(config.GetTmpDir(), "*.rdb")
	if err != nil {
		logger.Warnf("failed to create RDB file: %v", err)
		return err
	}
	defer file.Close()

	// todo: implement RDB generation logic
	cache := p.LoadAof()
	if cache == nil {
		return nil
	}

	encoder := rdbEncoder.NewEncoder(file).EnableCompress()
	err = encoder.WriteHeader()
	if err != nil {
		return err
	}
	auxMap := map[string]string{
		"redis-ver":    "6.0",
		"redis-bits":   "64",
		"aof-preamble": "0",
		"ctime":        strconv.FormatInt(time.Now().Unix(), 10),
	}
	if config.Properties.AofUseRdbPreamble {
		auxMap["aof-preamble"] = "1"
	}
	for k, v := range auxMap {
		err = encoder.WriteAux(k, v)
		if err != nil {
			return err
		}
	}
	cache.ForEach(func(key string, entity *kvcache.DataEntity, expiration *time.Time) bool {
		switch obj := entity.Data.(type) {
		case []byte:
			err = encoder.WriteStringObject(key, obj)
		}
		return err == nil
	})
	return nil
}
