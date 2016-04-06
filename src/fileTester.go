package main

import (
	// "./lib/colorprint"
	"./lib/fileshare"
	// "./lib/player"
	// "./lib/ui"
	// "./lib/utility"
	"os"
)

func main() {
	myAddr := os.Args[1]
	fileshare.FileSysStart(myAddr)
}
