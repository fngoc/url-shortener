package config

import "flag"

type flags struct {
	ServerAddress     string
	BaseResultAddress string
}

var Flags flags

func ParseArgs() {
	flag.StringVar(&Flags.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&Flags.BaseResultAddress, "b", "http://localhost:8080/", "base result server address")
	flag.Parse()
}
