package main

import (
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/server"
)

// main функция вызывается автоматически при запуске приложения
func main() {
	config.ParseArgs()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
