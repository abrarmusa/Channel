package main

import (
	"./lib/colorprint"
	// "./lib/player"
	"./lib/transfer"
	"net/http"
	// "./lib/ui"
	"./consts"
	"./lib/filemgmt"
	"bytes"
	// "./lib/utility"
	"fmt"
	"os"
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
			// go player.Run()
			for i := 1; i <= int(segNums); i++ {
				// fmt.Println("Getting seg ", i)
				vidSeg := transfer.GetVideoSegment("sample1.mp4", segNums, i, ":5000")
				// fmt.Println("got", vidSeg.Id)
				for j := 0; j < len(vidSeg.Body); j++ {
					// fmt.Println("Sending")
					player.ByteChan <- vidSeg.Body[j]
					// vid = append(vid, vidSeg.Body[j])
				}
				// fmt.Println("continuing")

				// // vidbytes = append(vidbytes, vidSeg.Body)
			}
			// go func(){

			// }
			// player.CloseStream <- 1
		}
		Run()
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// c := 0
	colorprint.Debug("Serving")
	w.Header().Set("Content-type", "video/mp4")
	var arr []byte
	// _, err := w.Write(vid)
	// utility.CheckError(err)
	b := bytes.NewBuffer([]byte{})

	for _, value := range vid {
		arr = append(arr, value)
		if len(arr) == consts.WindowSize {
			b.Write(arr)
			arr = []byte{}
			if _, err := b.WriteTo(w); err != nil { // <----- here!
				fmt.Fprintf(w, "%s", err)
			}
		} else {
			continue
		}

	}
	// _, err := w.Write(arr)
	// utility.CheckError(err)

}

// Run()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method starts the web player and the local streaming server
// -------------------
// INSTRUCTIONS:
// -------------------
// call player.Run()
func Run() {
	colorprint.Warning("Starting Player")
	// ByteChan = make(chan byte, consts.WindowSize)
	// CloseStream = make(chan int)
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8080", nil)
	// colorprint.Blue("The video " + Filename + " is streaming at http://localhost:8080/" + Filename)
	// switch runtime.GOOS {
	// case "linux":
	// 	err := exec.Command("xdg-open", "http://localhost:8080/").Start()
	// 	utility.CheckError(err)
	// case "windows", "darwin":
	// 	err := exec.Command("open", "http://localhost:8080/").Start()
	// 	utility.CheckError(err)
	// default:
	// 	colorprint.Alert("unsupported platform")
	// }

	// colorprint.Warning("Playing")
}
