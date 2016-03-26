package main

// package clientstream

import (
	// "net/rpc"
	// "os"
	// "errors"
	"fmt"
	"io/ioutil"
	"os"
	// "sync"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// HELPER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
}

func debug(str string) {
	fmt.Println(str)
}

func main() {
	file, err := os.Open("../samples/barbie.mp4")
	checkError(err)
	st, err := file.Stat()
	dat, err := ioutil.ReadFile("../samples/barbie.mp4")
	checkError(err)
	fmt.Println("Video has", len(dat), "segments. Size is", st.Sys())
}
