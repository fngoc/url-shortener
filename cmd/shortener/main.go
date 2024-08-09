package main

import (
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/server"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
	"github.com/fngoc/url-shortener/internal/logger"
)

// main функция вызывается автоматически при запуске приложения
func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	config.ParseArgs()

	if config.HasFlagAndEnvPostgresVariable() {
		if err := storage.InitializeDB(config.Flags.DBConf); err != nil {
			panic(err)
		}
	} else if config.HasFlagAndEnvFileVariable() {
		if err := storage.InitializeFileLocalStore(config.Flags.FilePath); err != nil {
			panic(err)
		}
	} else {
		if err := storage.InitializeInMemoryLocalStore(); err != nil {
			panic(err)
		}
	}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
