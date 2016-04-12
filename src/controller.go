package main

import (
	"./customChord"
	"./streamerClient"
	"./streamerServer"
	"os"
	"net/rpc"
	"log"
	"time"
	"fmt"
	"strconv"
)

// func generateStreamerQueue() {

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
	streamingClientAddress := "udp://127.0.0.1" + os.Args[4]
	name := os.Args[5]

	go customChord.Start(thisAddr, startNodeAddr, streamingServerAddress, streamingClientAddress)
	go streamerClient.ListenForStream(streamingClientAddress)
	go streamerServer.Start(streamingServerAddress, name)

	var shareFile string
	fmt.Println("====================================================")
	fmt.Println("Please enter the name of the file you wish to share:")
	fmt.Println("====================================================")
	fmt.Scanf("%s", &shareFile)

	// MOCK DONE: SPLIT THE FILE IN PARTS - returns array of segment filenames (?) - array of strings
	// TODO: DISTRIBUTE ALL PARTS USING CUSTOMCHORD
	if shareFile != "none" || shareFile != "" {
		totalSegments := streamerServer.GetFrames(shareFile)
		fmt.Printf("Total number of extracted frames from %s: %d\n", shareFile, totalSegments)
 // 	for _, seg := range segments {
 //        customChord.TransferFileSegment(seg)
 //    }
	}
	


	// for a node which holds several parts, just ask for the stream ONCE (?)

	var streamFile string
	fmt.Println("====================================================")
	fmt.Println("Please enter the name of the file you wish to stream:")
	fmt.Println("====================================================")
	fmt.Scanf("%s", &streamFile)
	const numParts = 2

    var handlers [2]*rpc.Client

    // Get streaming info for all parts
    for i := 0; i < numParts; i++ {
    	fn := streamFile + " " + strconv.FormatInt(int64(i), 10)
		addr := ""
		for addr == "" {
			log.Printf("Attempting to get stream server in 2 seconds...")
			time.Sleep(2 * time.Second)
			addr = customChord.GetStreamingServer(fn)
		}
		log.Println("Stream Server address: ", addr)

		//var handler *rpc.Client
		for handlers[i] == nil {
				log.Printf("Attempting to get rpc handler in 2 seconds...")
				time.Sleep(2 * time.Second)
				handlers[i] = streamerClient.GetRpcHandler(addr)
		}
	}

	// Start streaming
	// ASSUMPTION: EACH NODE STORES EITHER ATLEAST 300 SEQUENTIAL FRAMES, OR TILL THE END OF FRAME SEQUENCE
	for i := 0; i < numParts; i++ {
		streamerClient.StartStreaming(handlers[i], 0, strconv.FormatInt(int64(i*300), 10), streamingClientAddress)
	}



		// streamerClient.StartStreaming(handlers[0], 0, "0", streamingClientAddress)
		// streamerClient.StartStreaming(handlers[1], 0, "300", streamingClientAddress)
}