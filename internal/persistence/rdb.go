package persistence

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

const (
	// RDB 文件版本号
	RDB_VERSION = 1
	// 文件头部魔数
	RDB_MAGIC = "REDIS"
)

// RDB 文件格式：
// [MAGIC(5)][VERSION(4)][TIMESTAMP(8)][DB_SIZE(4)][ENTRIES...]
// ENTRY 格式：
// [KEY_LEN(4)][KEY][VALUE_TYPE(1)][VALUE_LEN(4)][VALUE]

// Save 保存数据库快照
func Save(db any, filename string) error {
	// 创建临时文件
	tmpFile := filename + ".tmp"
	file, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	// 写入文件头
	if err := writeHeader(file); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	// 获取数据库内容
	dbContent, err := json.Marshal(db)
	if err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	// 写入数据库大小
	if err := binary.Write(file, binary.LittleEndian, int32(len(dbContent))); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	// 写入数据库内容
	if _, err := file.Write(dbContent); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	// 确保所有数据都写入磁盘
	if err := file.Sync(); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return err
	}

	// 关闭文件
	file.Close()

	// 如果目标文件存在，先删除
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			os.Remove(tmpFile)
			return err
		}
	}

	// 原子性地重命名临时文件
	return os.Rename(tmpFile, filename)
}

// Load 加载数据库快照
func Load(filename string) (any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取并验证文件头
	if err := readHeader(file); err != nil {
		return nil, err
	}

	// 读取数据库大小
	var size int32
	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	// 读取数据库内容
	dbContent := make([]byte, size)
	if _, err := io.ReadFull(file, dbContent); err != nil {
		return nil, err
	}

	// 解析数据库内容
	var db any
	if err := json.Unmarshal(dbContent, &db); err != nil {
		return nil, err
	}

	return db, nil
}

// writeHeader 写入 RDB 文件头
func writeHeader(w io.Writer) error {
	// 写入魔数
	if _, err := w.Write([]byte(RDB_MAGIC)); err != nil {
		return err
	}

	// 写入版本号
	if err := binary.Write(w, binary.LittleEndian, int32(RDB_VERSION)); err != nil {
		return err
	}

	// 写入时间戳
	if err := binary.Write(w, binary.LittleEndian, time.Now().Unix()); err != nil {
		return err
	}

	return nil
}

// readHeader 读取并验证 RDB 文件头
func readHeader(r io.Reader) error {
	// 读取魔数
	magic := make([]byte, len(RDB_MAGIC))
	if _, err := io.ReadFull(r, magic); err != nil {
		return err
	}
	if string(magic) != RDB_MAGIC {
		return ErrInvalidRDBFormat
	}

	// 读取版本号
	var version int32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != RDB_VERSION {
		return ErrUnsupportedVersion
	}

	// 读取时间戳
	var timestamp int64
	if err := binary.Read(r, binary.LittleEndian, &timestamp); err != nil {
		return err
	}

	return nil
}

// 错误定义
var (
	ErrInvalidRDBFormat   = errors.New("invalid RDB format")
	ErrUnsupportedVersion = errors.New("unsupported RDB version")
)
