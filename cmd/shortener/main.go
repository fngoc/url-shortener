package main

import (
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/server"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
)

// main функция вызывается автоматически при запуске приложения
func main() {
	config.ParseArgs()

	if err := storage.InitializeLocalStore(config.Flags.FilePath); err != nil {
		panic(err)
	}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
