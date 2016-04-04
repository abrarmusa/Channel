package main

import (
	//"bytes"
	//"encoding/json"
	//"fmt"
	"log"
	"os/exec"
)

func getStream() {
	// ffplay udp://127.0.0.1:1234
	cmd := exec.Command("ffplay", "udp://127.0.0.1:1234")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting to get stream...")
	err = cmd.Wait()
	log.Printf("Getting stream finished with error: %v", err)
}

func main() {
	getStream()
}
