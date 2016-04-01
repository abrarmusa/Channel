package player

import (
	"../colorprint"
	"log"
	"net/http"
	"os"
	"time"
)

var Filepath string
var Filename string

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	video, err := os.Open(Filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer video.Close()
	http.ServeContent(w, r, "sample1.mp4", time.Now(), video)
}

func Run() {
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8080", nil)
	colorprint.Blue("Video " + Filename + " is available at http://localhost:8080/" + Filename)
}
