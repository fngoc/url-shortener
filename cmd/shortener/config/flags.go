package config

import (
	"flag"
	"os"
)

type flags struct {
	ServerAddress     string
	BaseResultAddress string
}

var Flags flags

func ParseArgs() {
	flag.StringVar(&Flags.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Flags.BaseResultAddress, "b", "http://localhost:8080", "base result server address")
	flag.Parse()

	serverAddressEnv, findAddress := os.LookupEnv("SERVER_ADDRESS")
	serverBaseUrlEnv, findBaseUrl := os.LookupEnv("BASE_URL")

	if findAddress {
		Flags.ServerAddress = serverAddressEnv
	}
	if findBaseUrl {
		Flags.BaseResultAddress = serverBaseUrlEnv
	}
}
