package main

import (
	"github.com/fngoc/url-shortener/cmd/shortener/server"
)

// main функция вызывается автоматически при запуске приложения
func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
