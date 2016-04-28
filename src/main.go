package main

import (
	"./lib/chordRPC"
	"./lib/filemgmt"
	"./lib/player"
	"./lib/transfer"
	//"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	chordAddress string
	ftAddress    string
	peerAddress  string
	//peerAddress1 	string
	vid []byte
)

func main() {

	vid = []byte{}

	if len(os.Args) < 4 {
		fmt.Printf("Usage : go run main.go <chordAddress> <ftAddress> <peerAddress>")
		os.Exit(-1)
	}

	chordAddress = os.Args[1]
	ftAddress = os.Args[2]
	peerAddress = os.Args[3]
	//peerAddress1 = os.Args[3]

	// Initialize local filesystem
	localFileSystem := transfer.Initialize(ftAddress, ":6666")
	filemgmt.ProcessLocalFiles(localFileSystem)
	filemgmt.PrintFileSysContents(localFileSystem)

	// Init chord
	go chordRPC.Start(chordAddress, peerAddress, ftAddress)

	var shareFile string
	fmt.Println("Please enter name of file you wish to share: ")
	fmt.Scan(&shareFile)
	if shareFile != "" {
		filemgmt.SplitFile(shareFile)
		fmt.Println("File splitting complete.")
	}

	// distribute the parts over all connected nodes
	fnArr := strings.Split(shareFile, ".")
	available, segNums, _ := transfer.CheckFileAvailability(shareFile, ftAddress)
	fmt.Printf("Total segments available to distribute: %d\n", segNums)
	// TODO: Delay
	if available {
		fmt.Println("Available. Commencing file load balancing...")

		// for all segs, distribute
		for i := 1; i <= int(segNums); i++ {
			filename := fnArr[0] + "_" + strconv.FormatInt(int64(i), 10)
			addr := chordRPC.GetAddressForSegment(filename)
			fmt.Println("Found node with address %s\n", addr)
			// now send file to addr
			if addr != ftAddress {
				//filemgmt.PrintFileSysContents(localFileSystem)
				vidSeg := transfer.GetVideoSegment(shareFile, segNums, i, ftAddress)
				transfer.SendVideoSegment(shareFile, addr, int(segNums), vidSeg)
				chordRPC.SaveToMap(filename, vidSeg.Body)
				fmt.Printf("Sent segment # %d\n", i)
			} else {
				fmt.Println("This node already stores the segment")
			}
		}
	}

	// stream a file
	var streamFile string
	fmt.Println("Please enter name of file you wish to stream: ")
	fmt.Scan(&streamFile)
	fnArr = strings.Split(streamFile, ".")
	if streamFile != "" {
		go player.Run()
		fmt.Printf("Preparing to stream %s\n", streamFile)

		for i := 1; i <= int(segNums); i++ {
			filename := fnArr[0] + "_" + strconv.FormatInt(int64(i), 10)
			addr := chordRPC.GetAddressForSegment(filename)

			// now get file segment from this node and push byte stream to vlc
			vidSeg := transfer.GetVideoSegment(streamFile, segNums, i, addr)
			for j := 0; j < len(vidSeg.Body); j++ {
				player.ByteChan <- vidSeg.Body[j]
				vid = append(vid, vidSeg.Body[j])
			}
		}
	}
	////////////////

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Node started. Enter some text to continue:")
	// text, _ := reader.ReadString('\n')
	// fmt.Println(text)

	// available, segNums, _ := transfer.CheckFileAvailability("sample1.mp4", peerAddress0)
	// fmt.Printf("Total segments available: %d\n", segNums)

	// if available {
	// 	fmt.Println("File is available")

	// 	go player.Run()

	// 	for i := 1; i <= int(segNums); i++ {
	// 		// get video segment from peer0
	// 		vidSeg := transfer.GetVideoSegment("sample1.mp4", segNums, i, peerAddress0)
	// 		// send video segment to peer1
	// 		transfer.SendVideoSegment("sample1.mp4", peerAddress1, int(segNums), vidSeg)

	// 		for j := 0; j < len(vidSeg.Body); j++ {
	// 			player.ByteChan <- vidSeg.Body[j]
	// 			vid = append(vid, vidSeg.Body[j])
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("File is unavailable")
	// }
	for {
	}
	fmt.Println("Exiting...")
}
