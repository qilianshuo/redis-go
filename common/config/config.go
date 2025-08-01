package config

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mirage208/redis-go/common/logger"
	"github.com/mirage208/redis-go/common/utils"
)

// ServerProperties defines global config properties
type ServerProperties struct {
	// for Public configuration
	RunID             string `cfg:"run-id"` // runID always different at every exec.
	Bind              string `cfg:"bind"`
	Port              int    `cfg:"port"`
	Dir               string `cfg:"dir"`
	AnnounceHost      string `cfg:"announce-host"`
	AppendOnly        bool   `cfg:"append-only"`
	AppendFilename    string `cfg:"append-filename"`
	AppendFsync       string `cfg:"append-fsync"`
	AofUseRdbPreamble bool   `cfg:"aof-use-rdb-preamble"`
	MaxClients        int    `cfg:"max-clients"`
	RequirePass       string `cfg:"require-pass"`
	Databases         int    `cfg:"databases"`
	RDBFilename       string `cfg:"db-filename"`
	ReplTimeout       int    `cfg:"repl-timeout"`

	// config file path
	CfPath string `cfg:"cf,omitempty"`
}

type ServerInfo struct {
	StartUpTime time.Time
}

// Properties holds global config properties
var Properties *ServerProperties
var EachTimeServerInfo *ServerInfo

func init() {
	// A few stats we don't want to reset: server startup time, and peak mem.
	EachTimeServerInfo = &ServerInfo{
		StartUpTime: time.Now(),
	}

	// default config
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
		RunID:      utils.RandString(40),
	}
}

func (p *ServerProperties) AnnounceAddress() string {
	if p.AnnounceHost != "" {
		return p.AnnounceHost + ":" + strconv.Itoa(p.Port)
	}
	return p.Bind + ":" + strconv.Itoa(p.Port)
}

func parse(src io.Reader) *ServerProperties {
	// TODO
	config := &ServerProperties{}

	// read config file
	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 { // separator found
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}
		if err := scanner.Err(); err != nil {
			logger.Fatal(err)
		}

		// parse format
		t := reflect.TypeOf(config)
		v := reflect.ValueOf(config)
		n := t.Elem().NumField()
		for i := 0; i < n; i++ {
			field := t.Elem().Field(i)
			fieldVal := v.Elem().Field(i)
			key, ok := field.Tag.Lookup("cfg")
			if !ok || strings.TrimLeft(key, " ") == "" {
				key = field.Name
			}
			value, ok := rawMap[strings.ToLower(key)]
			if ok {
				// fill config
				switch field.Type.Kind() {
				case reflect.String:
					fieldVal.SetString(value)
				case reflect.Int:
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err == nil {
						fieldVal.SetInt(intValue)
					}
				case reflect.Bool:
					boolValue := value == "yes"
					fieldVal.SetBool(boolValue)
				case reflect.Slice:
					if field.Type.Elem().Kind() == reflect.String {
						slice := strings.Split(value, ",")
						fieldVal.Set(reflect.ValueOf(slice))
					}
				default:
					logger.Fatal("unknown type")
				}
			}
		}
	}
	return config
}

// SetupConfig read config file and store properties into Properties
func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error(err)
		}
	}(file)
	Properties = parse(file)
	Properties.RunID = utils.RandString(40)
	configFilePath, err := filepath.Abs(configFilename)
	if err != nil {
		return
	}
	Properties.CfPath = configFilePath
	if Properties.Dir == "" {
		Properties.Dir = "."
	}
}
func GetTmpDir() string {
	return Properties.Dir + "/tmp"
}
