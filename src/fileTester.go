package main

import (
	// "./lib/colorprint"
	// "./lib/transfer"
	// "./lib/player"
	// "./lib/ui"
	// "./lib/utility"
	// "./consts"
	"./lib/filemgmt"
	// "fmt"
	// "os"
)

func main() {
	// myAddr := os.Args[1]
	// filename := "sample1.mp4"
	// nodeaddr := ":4000"
	// localFileSys := transfer.Initialize(myAddr)
	// avail, segNums, segsAvail := transfer.CheckFileAvailability(filename, nodeaddr)
	// var vidSegment utility.VidSegment
	// vidSegment = transfer.GetVideoSegment("sample.mp4", 45, ":3000")
	// transfer.Instr(myAddr)
	// fmt.Println(consts.DirPath + "/samples/sample1.mp4")
	// filemgmt.SplitFile(consts.DirPath + "/downloaded/sample1.mp4")
	filemgmt.ProcessLocalFiles()
	// filemgmt.
}
