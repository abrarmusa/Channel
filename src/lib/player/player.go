package player

import (
	"../../consts"
	"../colorprint"
	"../utility"
	// "io/ioutil"
	"bytes"
	"fmt"
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
	// waitingbytes, err := ioutil.ReadFile(consts.DirPath + "/waiting.mp4")
	// utility.CheckError(err)
	colorprint.Debug("Serving")
	w.Header().Set("Content-type", "video/mp4")
	var arr []byte = []byte{}
	// _, err := w.Write(vid)
	// utility.CheckError(err)
	b := bytes.NewBuffer([]byte{})
	// b.Write(waitingbytes)
	// _, err = b.WriteTo(w)
	// colorprint.Debug("RRRdd")
	// utility.CheckError(err)
	c := 0
	for tmp := range ByteChan {
		arr = append(arr, tmp)
		if len(arr) == 4096 {
			b.Write(arr)
			arr = []byte{}
			_, err := b.WriteTo(w)
			fmt.Printf("\r%s", r.Header, "\nrSending 4096 bytes. Sequence number is %d", c)
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
	// colorprint.Debug("------------------------------------------------------------------")
	// for {
	// 	select {
	// 	case tmp := <-ByteChan:
	// 		// colorprint.Debug("RECEIVED")
	// 		c++
	// 		arr = append(arr, tmp)
	// 		if len(arr) == 4096 {
	// 			b.Write(arr)
	// 			arr = []byte{}
	// 			_, err := b.WriteTo(w)
	// 			fmt.Printf("%s", r.Header)
	// 			fmt.Printf("\rSending 4096 bytes. Sequence number is %d", c)
	// 			if err != nil {
	// 				colorprint.Debug("err")
	// 			}
	// 			// utility.CheckError(err)
	// 		}
	// 	case <-CloseStream:
	// 		for {
	// 			select {
	// 			case tmp2 := <-ByteChan:
	// 				arr = append(arr, tmp2)
	// 				if len(arr) == 4096 {
	// 					b.Write(arr)
	// 					arr = []byte{}
	// 					_, err := b.WriteTo(w)
	// 					fmt.Printf("%s", r.Header)
	// 					// fmt.Printf("\rSending 4096 bytes. Sequence number is %d", c)
	// 					if err != nil {
	// 						colorprint.Debug("err")
	// 					}
	// 					// utility.CheckError(err)
	// 				}
	// 			default:
	// 				break
	// 			}
	// 		}
	// 		colorprint.Debug("------------------------------------------------------------------")
	// 		b.Write(arr)
	// 		arr = []byte{}
	// 		_, err := b.WriteTo(w)
	// 		utility.CheckError(err)
	// 		colorprint.Debug("VIDEO STREAM COMPLETED. CLOSING STREAM")
	// 		w.Write([]byte("Video Completed"))
	// 		colorprint.Debug("------------------------------------------------------------------")
	// 		break
	// 	default:
	// 		// b.Write(waitingbytes)
	// 		// _, err := b.WriteTo(w)
	// 		// colorprint.Debug("RRRdd")
	// 		// utility.CheckError(err)
	// 		// continue
	// 		fmt.Printf("%s", r.Header)
	// 		fmt.Printf("\rWaiting for stream")
	// time.Sleep(300 * time.Millisecond)
	// 		continue
	// 	}
	// }
	colorprint.Debug("CLOSING")
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
	ByteChan = make(chan byte, consts.WindowSize*100000)
	CloseStream = make(chan int)
	http.HandleFunc("/", ServeHTTP)
	go http.ListenAndServe(":8080", nil)
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
