package main

import (
	// "./lib/colorprint"
	"./lib/fileshare"
	// "./lib/player"
	// "./lib/ui"
	"./lib/utility"
	"os"
)

func main() {
	myAddr := os.Args[1]
	filename := "sample1.mp4"
	nodeaddr := ":4000"
	fileshare.FileSysStart(myAddr)
	avail, segNums, segsAvail := fileshare.CheckFileAvailability(filename, nodeaddr)
	var vidSegment utility.VidSegment
	vidSegment = fileshare.GetVideoSegment("sample.mp4", 45, ":3000")
}
