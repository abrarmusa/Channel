package main
/*
  UBC CS416 Distributed Systems Project Source Code

  @author: Abrar, Ito, Mimi, Shariq
  @date: Mar. 1 2016 - Apr. 11 2016.

  Usage:
    go run stream_client.go [starter-node ip:port]"

    [starter-node ip:port] : the entry point node's ip/port combo

  Copy/paste for quick testing:
    "go run stream_client.go :6666" <- connect to node listening at :6666
*/

import (
	"clientstream"
	"fmt"
  "lib/fileshare"
)

func main() {
  if len(os.Args) != 3 {
    fmt.Println("Usage: go run stream_client.go [starter-node ip:port]")
    os.Exit(-1)
  } else {
    startAddr := os.Args[1]                   // ip:port of initial node

  }
}
