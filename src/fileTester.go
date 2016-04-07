package main

import (
	// "./lib/colorprint"
	"./lib/fileshare"
	// "./lib/player"
	// "./lib/ui"
	// "./lib/utility"
	"os"
	// "./consts"
	// "io/ioutil"
)

func main() {
	// filename := "1319.jpg"
	// nodeaddr := "192.168.0.49:4000"
	myAddr := os.Args[1]
	// str := consts.DirPath + "/downloaded/" +filename
	fileshare.FileSysStart(myAddr)
	// _, segNums, _ := fileshare.CheckFileAvailability(filename, nodeaddr)
	// var vidSegment utility.VidSegment
	// var b []byte
	// for i:=0; i < int(segNums); i++ {
	// 	vidSegment = fileshare.GetVideoSegment(filename, i, nodeaddr)
	// 	for j:=0; j < len(vidSegment.Body); j++ {
	// 		b = append(b, vidSegment.Body[j])
	// 	}
	// }
	// err := ioutil.WriteFile(str, b, 0777)
	// utility.CheckError(err)
	fileshare.Instr(myAddr)
	
}
