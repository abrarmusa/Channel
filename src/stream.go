// REQUIRES A LOT OF REFACTORIZATION OFCOURSE
// NEED TO MAKE FUNCTIONS MORE MODULAR - WITH PARAMS FOR DIFFERENT FILENAMES AND OUTPUT FILENAMES


package main

import (
	"log"
	"os/exec"
)

func getFrames() {
	// ffmpeg -i sample.mp4 -r 100 -f image2 output/%05d.png
	cmd := exec.Command("ffmpeg", "-i", "FFMPEG/sample.mp4", "-r", "100", "-f",
		"image2", "FFMPEG/output/%05d.png")
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting for video to finish processing into individual frames...")
	err = cmd.Wait()
	log.Printf("Frame processing finished with error: %v", err)
}

func startStream() {
	// ffmpeg -re -i FFMPEG/output/output%05d.png -r 10 -vcodec mpeg4 -f mpegts udp://127.0.0.1:1234
	cmd := exec.Command("ffmpeg", "-re", "-i", "FFMPEG/output/%05d.png", "-r", "10", 
		"-vcodec", "mpeg4", "-f", "mpegts",  "udp://127.0.0.1:1234")
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting to start streaming frames...")
	err = cmd.Wait()
	log.Printf("Frame streaming finished with error: %v", err)
	
}

func main() {
	getFrames()
	startStream()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}