package main

import (
	// "./lib/colorprint"
	"./lib/player"
	"./lib/transfer"
	// "./lib/ui"
	// "./lib/utility"
	// "./consts"
	"./lib/filemgmt"
	"fmt"
	"os"
	// "./consts"
	// "io/ioutil"
)

func noder(myAddr string) {
	// filename := "sample1.mp4"
	// nodeaddr := ":4000"
	localFileSys := transfer.Initialize(myAddr)
	// avail, segNums, segsAvail := transfer.CheckFileAvailability(filename, nodeaddr)
	// var vidSegment utility.VidSegment
	// vidSegment = transfer.GetVideoSegment("sample.mp4", 45, ":3000")
	// transfer.Instr(myAddr)
	// fmt.Println(consts.DirPath + "/samples/sample1.mp4")
	// filemgmt.SplitFile(consts.DirPath + "/downloaded/sample1.mp4")
	filemgmt.ProcessLocalFiles(localFileSys)
	filemgmt.PrintFileSysContents(localFileSys)
	transfer.Instr()
}
func main() {
	// filename := "1319.jpg"
	// nodeaddr := "192.168.0.49:4000"
	myAddr := os.Args[1]

	if os.Args[2] == "2" {
		noder(myAddr)
	} else {

		// filename := "sample1.mp4"
		// nodeaddr := ":4000"
		localFileSys := transfer.Initialize(myAddr)
		// avail, segNums, segsAvail := transfer.CheckFileAvailability(filename, nodeaddr)
		// var vidSegment utility.VidSegment
		// vidSegment = transfer.GetVideoSegment("sample.mp4", 45, ":3000")
		// transfer.Instr(myAddr)
		// fmt.Println(consts.DirPath + "/samples/sample1.mp4")
		// filemgmt.SplitFile(consts.DirPath + "/downloaded/sample1.mp4")
		filemgmt.ProcessLocalFiles(localFileSys)
		filemgmt.PrintFileSysContents(localFileSys)
		// transfer.Instr()
		avail, segNums, _ := transfer.CheckFileAvailability("sample1.mp4", ":5000")
		fmt.Println("SEGMENTS AVAIL: ", segNums)
		if avail {
			// now get the file segments from the node 5000

			// var vidbytes []byte
			go player.Run()
			for i := 1; i <= int(segNums); i++ {
				// fmt.Println("Getting seg ", i)
				vidSeg := transfer.GetVideoSegment("sample1.mp4", segNums, i, ":5000")
				// fmt.Println("got", vidSeg.Id)
				for j := 0; j < len(vidSeg.Body); j++ {
					// fmt.Println("Sending")
					player.ByteChan <- vidSeg.Body[j]
				}
				// fmt.Println("continuing")

				// // vidbytes = append(vidbytes, vidSeg.Body)
			}
			// go func(){

			// }
			player.CloseStream <- 1
		}
	}
}