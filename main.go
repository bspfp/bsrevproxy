package main

import (
	"bspfp/bsrevproxy/internal/config"
	"bspfp/bsrevproxy/internal/server"
	"bspfp/gosimplelog"
	"fmt"
	"os"
)

func main() {
	logCloser, err := gosimplelog.InitLogFile("./log", "bsrevproxy.log", 10)
	if err != nil {
		fmt.Println("failed to init log file:", err)
		os.Exit(1)
	}
	defer logCloser.Close()

	config.ParseFlag()

	server.StartServer()
}
