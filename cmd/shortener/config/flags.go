package config

import (
	"flag"
	"os"
)

type flags struct {
	ServerAddress     string
	BaseResultAddress string
	FilePath          string
	DBConf            string
}

var Flags flags

func ParseArgs() {
	flag.StringVar(&Flags.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Flags.BaseResultAddress, "b", "http://localhost:8080", "base result server address")
	flag.StringVar(&Flags.FilePath, "f", "data.json", "file path")
	flag.StringVar(&Flags.DBConf, "d", "host=localhost user=postgres password=postgres dbname=test_db sslmode=disable", "db params")
	flag.Parse()

	serverAddressEnv, findAddress := os.LookupEnv("SERVER_ADDRESS")
	serverBaseURLEnv, findBaseURL := os.LookupEnv("BASE_URL")
	filePathEnv, findFilePath := os.LookupEnv("FILE_STORAGE_PATH")
	DBEnv, findDBConf := os.LookupEnv("DATABASE_DSN")

	if findAddress {
		Flags.ServerAddress = serverAddressEnv
	}
	if findBaseURL {
		Flags.BaseResultAddress = serverBaseURLEnv
	}
	if findFilePath {
		Flags.FilePath = filePathEnv
	}
	if findDBConf {
		Flags.DBConf = DBEnv
	}
}
