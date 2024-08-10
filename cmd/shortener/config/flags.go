package config

import (
	"flag"
	"github.com/fngoc/url-shortener/internal/logger"
	"os"
)

type flags struct {
	ServerAddress     string
	BaseResultAddress string
	FilePath          string
	DBConf            string
}

var Flags flags

const defaultPostgresParams string = "host=localhost user=postgres password=postgres dbname=test_db sslmode=disable"
const defaultFileParams string = "data.json"

func ParseArgs() {
	flag.StringVar(&Flags.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Flags.BaseResultAddress, "b", "http://localhost:8080", "base result server address")
	flag.StringVar(&Flags.FilePath, "f", defaultFileParams, "file path")
	flag.StringVar(&Flags.DBConf, "d", defaultPostgresParams, "db params")
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
	logger.Log.Info("Parse argument's is done")
}

func HasFlagOrEnvPostgresVariable() bool {
	_, find := os.LookupEnv("DATABASE_DSN")
	if Flags.DBConf != defaultPostgresParams || find {
		return true
	}
	return false
}

func HasFlagOrEnvFileVariable() bool {
	_, find := os.LookupEnv("FILE_STORAGE_PATH")
	if Flags.FilePath != defaultFileParams || find {
		return true
	}
	return false
}
