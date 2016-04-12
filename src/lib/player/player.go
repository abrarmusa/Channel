package player

import (
	"../../consts"
	"../colorprint"
	"../utility"
	"fmt"
	"net/http"
	// "os/exec"
	// "runtime"
	"time"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// GLOBAL VARS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
var Filepath string
var Filename string
var ByteChan chan byte
var CloseStream chan int

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// VIDEO STREAMING METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// ServeHTTP(w http.ResponseWriter, r *http.Request)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method takes in incoming byte from the ByteChan channel and streams those video bytes to the webpage at localhost:8080/
// -------------------
// INSTRUCTIONS:
// -------------------
// NOTE: You must call player.Run() before calling this method.
// Call player.Run(). Then pass in your byte into player.ByteChan
//
// E.g code:
// tmp := make(byte, 5, 5)
// player.ByteChan <- tmp
//
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := 0
	w.Header().Set("Content-type", "video/mp4")
	var arr []byte
	// var arr2 []byte
	var prev []byte
	for {
		select {
		case tmp := <-ByteChan:
			arr = append(arr, tmp)
			c++
			if c == consts.WindowSize {
				t := time.Now().String()
				colorprint.Debug("------------------------------------------------------------------")
				colorprint.Debug(">> " + t + " <<")
				fmt.Println("VIDEOSTREAM: Sending ", consts.WindowSize, " bytes of the video.")
				_, err := w.Write(arr)

				// fix error here
				fmt.Println("IS IT BROKEN?", len(arr), " - ", c)
				utility.CheckError(err)
				fmt.Println("NOPE")
				prev = arr
				arr = []byte{}
				c = 0
				colorprint.Debug("------------------------------------------------------------------")
				// time.Sleep(5000 * time.Millisecond)
			} else {
				w.Write(prev)
			}
		case <-CloseStream:
			colorprint.Debug("------------------------------------------------------------------")
			_, err := w.Write(arr)
			utility.CheckError(err)
			colorprint.Debug("VIDEO STREAM COMPLETED. CLOSING STREAM")
			w.Write([]byte("Video Completed"))
			colorprint.Debug("------------------------------------------------------------------")
			break
		default:
			_, err := w.Write(prev)
			// fix error here
			fmt.Println("DO THIS", len(arr))
			utility.CheckError(err)
		}
	}
	close(ByteChan)
	close(CloseStream)

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
	ByteChan = make(chan byte, consts.WindowSize)
	CloseStream = make(chan int)
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
