package player

import (
	"../../consts"
	"../colorprint"
	"../utility"
	// "io/ioutil"
	"bytes"
	// "fmt"
	"net/http"
	// "time"
	// "os/exec"
	// "runtime"
	// "time"
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
	colorprint.Debug("Serving")
	w.Header().Set("Content-type", "video/mp4")
	var arr []byte = []byte{}
	b := bytes.NewBuffer([]byte{})
	for tmp := range ByteChan {
		arr = append(arr, tmp)
		if len(arr) == 4096 {
			b.Write(arr)
			arr = []byte{}
			_, err := b.WriteTo(w)
			if err != nil {
				colorprint.Debug("err")
			}
			utility.CheckError(err)
		}
	}
	colorprint.Debug("------------------------------------------------------------------")
	b.Write(arr)
	arr = []byte{}
	_, err := b.WriteTo(w)
	utility.CheckError(err)
	colorprint.Debug("VIDEO STREAM COMPLETED. CLOSING STREAM")
	w.Write([]byte("Video Completed"))
	colorprint.Debug("CLOSING")
	close(ByteChan)
	close(CloseStream)
	for {}

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
	ByteChan = make(chan byte, consts.WindowSize*100000)
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
