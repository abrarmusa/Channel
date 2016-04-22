package main

import (
	"./lib/customChord"
	"./lib/streamerClient"
	"./lib/streamerServer"
	//"./lib/transfer"
	//"./lib/utility"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
	//"runtime"
)

var name string
var ftAddr string
var streamingServerAddress string

type VidFrames struct {
	Name        string
	FrameStart  string
	TotalFrames int64
	Data        []byte
}

func frameToBytesArray(folderName string) map[string][]byte {
	frameBytesMap := make(map[string][]byte)

	folderPath := "FFMPEG/NodesData/" + name + "/" + folderName
	files, err := ioutil.ReadDir(folderPath)
	checkError(err)
	var frameBytes []byte
	for i, file := range files {
		filename := file.Name()
		fmt.Printf("File at index %d: %s\n", i, filename)
		frameBytes, err = ioutil.ReadFile(folderPath + "/" + filename)
		checkError(err)
		frameBytesMap[filename] = frameBytes
	}
	return frameBytesMap
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
	1. my upd address
	2. starter node udp address

	3. udp address of where im going to be listening for udp streams
	4. my streamerServer address

	5. node name used for streamer server
*/
func main() {

	//runtime.GOMAXPROCS(4)

	thisAddr := os.Args[1]
	startNodeAddr := os.Args[2]
	streamingServerAddress = os.Args[3]
	streamingClientAddress := "udp://127.0.0.1" + os.Args[4]
	name = os.Args[5]
	//ftAddr = os.Args[6]

	//_ = transfer.Initialize(ftAddr, name)

	go customChord.Start(thisAddr, startNodeAddr, streamingServerAddress, streamingClientAddress, ftAddr, name)
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

		// Calculate number of nodes to distribute on. Each node can hold roughly 200 frames
		totalNodes := totalSegments / 200
		fmt.Println("Total number of nodes this file will be distributed on: ", totalNodes)
		fnArr := strings.Split(shareFile, ".")

		// Get all the byte arrays for this file
		frameBytesArr := frameToBytesArray(fnArr[0])

		var fileNodes = make([]string, totalNodes)

		for i := 0; i < int(totalNodes); i++ {
			filenameWithNodeSegment := fnArr[0] + " " + strconv.FormatInt(int64(i), 10)
			fileNodes[i] = customChord.GetTransferFileSegmentAddr(filenameWithNodeSegment)
			fmt.Println("Address of file node to transfer: ", fileNodes[i])

			// transfer 200 frames
			// filename := fmt.Sprintf("%05d.png", int64(i))
			//seqNum := fmt.Sprintf("%d", int64(i+1))
			fn := fnArr[0] + " " + strconv.FormatInt(int64(i), 10)
			var addr string

			for addr == "" {
				log.Printf("Attempting to get ft server in 2 seconds...")
				time.Sleep(2 * time.Second)
				addr = customChord.GetStreamingServer(fn)
			}
			fmt.Printf("File %s to be transferred to %s\n", fn, addr)

			if addr != streamingServerAddress {
				var handler *rpc.Client
				for handler == nil {
					handler = streamerClient.GetRpcHandler(addr)
					if handler != nil {
						break
					}
					log.Printf("Attempting to get rpc handler for ft in 2 seconds...")
					time.Sleep(2 * time.Second)
				}
				var limit int
				if limit = i*200 + 201; i == int(totalNodes)-1 {
					limit = int(totalSegments) + 1
				}
				for j := i*200 + 1; j < limit; j++ {
					filename := fmt.Sprintf("%05d.png", int64(j))
					folderFilePath := fnArr[0] + " " + filename
					//fmt.Println("Saving to server...")
					streamerClient.SaveToServer(handler, name, folderFilePath, frameBytesArr[filename], ftAddr)
					//customChord.SaveToStore(fnArr[0], strconv.FormatInt(int64(i), 10), frameBytesArr[filename])
					//fmt.Println("Saved!!!!!!!!!!!!!!!!!!!!!!")
				}
			} else {
				fmt.Println("Already saved in this node")
			}
		}
	}

	// for a node which holds several parts, just ask for the stream ONCE (?)

	var streamFile string
	fmt.Println("====================================================")
	fmt.Println("Please enter the name of the file you wish to stream:")
	fmt.Println("====================================================")
	fmt.Scanf("%s", &streamFile)
	const numParts = 2

	var handlers [2]*rpc.Client
	fnArr := strings.Split(streamFile, ".")
	// Get streaming info for all parts
	for i := 0; i < numParts; i++ {
		fn := fnArr[0] + " " + strconv.FormatInt(int64(i), 10)
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
		streamerClient.StartStreaming(handlers[i], streamFile, 0, strconv.FormatInt(int64(i*300), 10), streamingClientAddress)
	}
	// streamerClient.StartStreaming(handlers[0], 0, "0", streamingClientAddress)
	// streamerClient.StartStreaming(handlers[1], 0, "300", streamingClientAddress)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
