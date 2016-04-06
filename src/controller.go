package main

import (
	"./customChord"
	"./streamerClient"
	"./streamerServer"
	"os"
	"net/rpc"
	"log"
	"time"
)

// streamerClient.ListenForStream("udp://127.0.0.1:1234")
// 	handler0 := streamerClient.GetRpcHandler(":1342")
// 	handler1 := streamerClient.GetRpcHandler(":1354")
// 	streamerClient.StartStreaming(handler0, 0, "0")
// 	streamerClient.StartStreaming(handler1, 1, "300")

// This method responds to the user input requests
// func serveUser(nodeRPC string, nodeUDP string) {
// 	var input string
// 	var fname string
// 	for {
// 		fmt.Println()
// 		colorprint.Info(">>>> Please enter a command")
// 		input, err := fmt.Scanln(&input)
// 		//fmt.Scan(&input)
// 		cmd := input[0]

// 		switch input {
// 			case "stream":
// 				filename := input[1]
// 				addr := customChord.GetStreamingServer(filename)
// 				handler := streamerClient.GetRpcHandler(addr)
// 				streamerClient.StartStreaming(handler, 0, "0")
// 		}

// 	}
// }
/*
	1. my upd address
	2. starter node udp address

	3. udp address of where im going to be listening for udp streams
	4. my streamerServer address

	5. node name used for streamer server
*/
func main() {
	thisAddr := os.Args[1]
	startNodeAddr := os.Args[2]
	streamingServerAddress := os.Args[3]
	streamingClientAddress := os.Args[4]
	name := os.Args[5]

	go customChord.Start(thisAddr, startNodeAddr, streamingServerAddress, streamingClientAddress)
	go streamerClient.ListenForStream("udp://127.0.0.1:1234")
	go streamerServer.Start(streamingServerAddress, name)

	addr := ""
	for addr == "" {
		log.Printf("Attempting to get stream server in 2 seconds...")
		time.Sleep(2 * time.Second)
		addr = customChord.GetStreamingServer("0")
	}

	if addr != "" {
		var handler *rpc.Client
		for handler == nil {
			log.Printf("Attempting to get stream server in 2 seconds...")
			time.Sleep(2 * time.Second)
			handler = streamerClient.GetRpcHandler(addr)
		}
		streamerClient.StartStreaming(handler, 0, "0", streamingClientAddress)
	} else {
		log.Printf("GetStreamingServer returned empty address. Quitting.")
	}
}