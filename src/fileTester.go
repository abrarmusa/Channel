package main

import (
	// "./lib/colorprint"
	// "./lib/player"
	"./lib/transfer"
	// "net/http"
	// "./lib/ui"
	"./lib/filemgmt"
	// "bytes"

	// "./lib/utility"
	// "fmt"
	"os"
	// "io/ioutil"
)

var vid []byte

func noder(myAddr string) {
	// filename := "sample1.mp4"
	// nodeaddr := ":4000"
	localFileSys := transfer.Initialize(myAddr, "node1")
	// avail, segNums, segsAvail := transfer.CheckFileAvailability(filename, nodeaddr)
	// var vidSegment utility.VidSegment
	// vidSegment = transfer.GetVideoSegment("sample.mp4", 45, ":3000")
	// transfer.Instr(myAddr)
	// fmt.Println(consts.DirPath + "/samples/sample1.mp4")
	filemgmt.SplitFile("sample1.mp4")
	filemgmt.ProcessLocalFiles(localFileSys)
	filemgmt.PrintFileSysContents(localFileSys)
	transfer.Instr()
}

func main() {
	vid = []byte{}
	myAddr := os.Args[1]
	// peerAddr := os.Args[2]s

	if os.Args[2] == "2" {
		noder(myAddr)
	}
	// }
	// else {
	// 	localFileSys := transfer.Initialize(myAddr, ":1234")
	// 	filemgmt.ProcessLocalFiles(localFileSys)
	// 	fmt.Println("HEKK")
	// 	filemgmt.PrintFileSysContents(localFileSys)
	// 	fmt.Println("HEKK")
	// 	// transfer.Instr()
	// 	avail, segNums, _ := transfer.CheckFileAvailability("sample1.mp4", ":5000")
	// 	fmt.Println("SEGMENTS AVAIL:", segNums)
	// 	if avail {
	// 		// Ready the streamer in a separate goroutine
	// 		go player.Run()
	// 		// now get the file segments from the node 5000
	// 		for i := 1; i <= int(segNums); i++ {
	// 			// fmt.Println("Getting seg ", i)
	// 			vidSeg := transfer.GetVideoSegment("sample1.mp4", segNums, i, ":5000")
	// 			transfer.SendVideoSegment("sample.mp4", peerAddr, int(segNums), vidSeg)
	// 			// transfer.SendVideoSegment(fname string, nodeAdd string, segNums int, segment utility.VidSegment)

	// 			// fmt.Println("got", vidSeg.Id)
	// 			for j := 0; j < len(vidSeg.Body); j++ {
	// 				// fmt.Println("Sending")
	// 				player.ByteChan <- vidSeg.Body[j]
	// 				vid = append(vid, vidSeg.Body[j])
	// 			}
	// 		}
	// 		fmt.Println("CLOSING")
	// 		// close(player.ByteChan)
	// 	}
	// //}

}
