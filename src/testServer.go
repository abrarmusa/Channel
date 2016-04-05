package main

import (
	"./streamerServer"
	"os"
	"fmt"
)

func main() {

	if len(os.Args) != 3 {
	    fmt.Println("Usage: go run testServer.go [node ip:port] [node name]")
	    os.Exit(-1)
  	}

	streamerServer.Start(os.Args[1], os.Args[2])
}