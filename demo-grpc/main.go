package main

import (
	"flag"

	"github.com/gohive/demo-grpc/cmd"
)

func main() {
	configPath := flag.String("config", "./config.toml", "path to config file")
	flag.Parse()

	cmd.Run(*configPath)
}
