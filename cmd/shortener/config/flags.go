package config

import (
	"flag"
	"os"
)

type flags struct {
	ServerAddress     string
	BaseResultAddress string
	FilePath          string
}

var Flags flags

func ParseArgs() {
	flag.StringVar(&Flags.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Flags.BaseResultAddress, "b", "http://localhost:8080", "base result server address")
	flag.StringVar(&Flags.FilePath, "f", "data.json", "file path")
	flag.Parse()

	serverAddressEnv, findAddress := os.LookupEnv("SERVER_ADDRESS")
	serverBaseURLEnv, findBaseURL := os.LookupEnv("BASE_URL")
	filePathEnv, findFilePath := os.LookupEnv("FILE_STORAGE_PATH")

	if findAddress {
		Flags.ServerAddress = serverAddressEnv
	}
	if findBaseURL {
		Flags.BaseResultAddress = serverBaseURLEnv
	}
	if findFilePath {
		Flags.FilePath = filePathEnv
	}
}
