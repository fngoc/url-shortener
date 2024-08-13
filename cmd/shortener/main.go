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

	if config.HasFlagOrEnvPostgresVariable() {
		if err := storage.InitializeDB(config.Flags.DBConf); err != nil {
			logger.Log.Fatal(err.Error())
		}
	} else if config.HasFlagOrEnvFileVariable() {
		if err := storage.InitializeFileLocalStore(config.Flags.FilePath); err != nil {
			logger.Log.Fatal(err.Error())
		}
	} else {
		if err := storage.InitializeInMemoryLocalStore(); err != nil {
			logger.Log.Fatal(err.Error())
		}
	}
	if err := server.Run(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
