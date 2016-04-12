package main

import (
	// "./lib/colorprint"
	"./lib/player"
	"./lib/transfer"
	// "net/http"
	// "./lib/ui"
<<<<<<< HEAD
	// "./consts"
	"./lib/filemgmt"
	// "bytes"

	// "./lib/utility"
	"fmt"
=======
	// "./lib/utility"
>>>>>>> bbdc7b9b739a41183df9a242fe06ee97401ca755
	"os"
	// "./consts"
	// "io/ioutil"
)

var vid []byte

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
<<<<<<< HEAD

	vid = []byte{}
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
		fmt.Println("SEGMENTS AVAIL:", segNums)
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
					vid = append(vid, vidSeg.Body[j])
				}
				// fmt.Println("continuing")

			}
			// go func(){

			// }
			fmt.Println("CLOSING")
			close(player.ByteChan)
		}
		// Run()
	}

}
