package main

import (
	"./lib/filemgmt"
	"./lib/player"
	"./lib/transfer"
	"bufio"
	"fmt"
	"os"
)

var (
	nodeAddress  string
	peerAddress0 string
	peerAddress1 string
	vid          []byte
)

func main() {

	vid = []byte{}

	if len(os.Args) < 4 {
		fmt.Printf("Usage : go run main.go <nodeAddress> <peerAddress0>")
		os.Exit(-1)
	}

	nodeAddress := os.Args[1]
	peerAddress0 := os.Args[2]
	peerAddress1 := os.Args[3]

	// Initialize local filesystem
	localFileSystem := transfer.Initialize(nodeAddress, ":6666")
	filemgmt.ProcessLocalFiles(localFileSystem)
	filemgmt.PrintFileSysContents(localFileSystem)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Node started. Enter some text to continue:")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	available, segNums, _ := transfer.CheckFileAvailability("sample1.mp4", peerAddress0)
	fmt.Printf("Total segments available: %d\n", segNums)

	if available {
		fmt.Println("File is available")

		go player.Run()

		for i := 1; i <= int(segNums); i++ {
			// get video segment from peer0
			vidSeg := transfer.GetVideoSegment("sample1.mp4", segNums, i, peerAddress0)
			// send video segment to peer1
			transfer.SendVideoSegment("sample1.mp4", peerAddress1, int(segNums), vidSeg)

			for j := 0; j < len(vidSeg.Body); j++ {
				player.ByteChan <- vidSeg.Body[j]
				vid = append(vid, vidSeg.Body[j])
			}
		}
	} else {
		fmt.Println("File is unavailable")
	}

	fmt.Println("Exiting...")
}
