package main

import (
	"fmt"
	"os"

	"github.com/qilianshuo/redis-go/common/logger"
	"github.com/qilianshuo/redis-go/common/utils"
	"github.com/qilianshuo/redis-go/internal/config"
	"github.com/qilianshuo/redis-go/internal/transport"
)

var banner = `
██████╗ ███████╗██████╗ ██╗███████╗       ██████╗  ██████╗ 
██╔══██╗██╔════╝██╔══██╗██║██╔════╝      ██╔════╝ ██╔═══██╗
██████╔╝█████╗  ██║  ██║██║███████╗█████╗██║  ███╗██║   ██║
██╔══██╗██╔══╝  ██║  ██║██║╚════██║╚════╝██║   ██║██║   ██║
██║  ██║███████╗██████╔╝██║███████║      ╚██████╔╝╚██████╔╝
╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝       ╚═════╝  ╚═════╝
`

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           6389,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
	RunID:          utils.RandString(40),
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "redis-go",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	// Load configuration
	// Check if CONFIG environment variable is set
	// If not set, check for redis.conf in the current directory
	// If redis.conf exists, load it; otherwise, use default properties
	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {
		if utils.FileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFilename)
	}

	// Start the server
	err := transport.ListenAndServeWithSignal(&transport.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, transport.NewHandler())
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}
