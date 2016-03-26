package main

// package clientstream

import (
	"net/rpc"
	// "os"
	"errors"
	"fmt"
	"sync"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//  STRUCTS & TYPES
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// --> VidSegment <---
// -------------------
// DESCRIPTION:
// -------------------
// This struct holds a particular part of a video file. Id refers to the segment id and the body refers to the actual video bytes
type VidSegment struct {
	Id   int
	Body []byte
}

// --> Video <---
// --------------
// DESCRIPTION:
// -------------------
// This struct holds the VidSegments of a particular video file. SegNums refers to the total number of segments in the entire video.
// This info is used to reorder the segments when playing the video stream.
type Video struct {
	Name      string
	SegNums   int64
	SegsAvail map[int]int64
	Segments  map[int]VidSegment
}

// --> FileSys <---
// ----------------
// DESCRIPTION
// -------------------
// This struct represents the local FileSystem to hold the Video's. Each FileSys object has an id and a map of Files with the keys to the map being the
// actual filename.
// This info is used to check if a node actually has the file
type FileSys struct {
	Id    int
	Files map[string]Video
}

// --> Response <---
// -----------------
// DESCRIPTION:
// -------------------
// This struct represents the response that an RPC call will write to. It is used to check if a node has a particular file and if it does, which parts
// of that file it has in its local filesystem.
type Response struct {
	Avail     bool
	SegNums   int64
	SegsAvail map[int]int64
}

// --> Service <---
// ----------------
// DESCRIPTION:
// -------------------
// This type just holds an integer to use for registering the RPC Service
type Service int

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// GLOBAL VARS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

var localFileSys FileSys
var fileSysLock *sync.RWMutex

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// INBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// (service *Service) localFileAvailability(filename string, response *Response) error
// -----------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to an rpc Call for a particular segment of a file. It first looks up and checks if the video file is available. If
// the file is available, it continues on to see if the segment is available. If the segment is available, it returns a response with the VidSegment.
// In case of unavailability, it will either return an error saying "File Unavailable." or "Segment unavailable." depending on what was unavailable.
// The method locks the local filesystem for all Reads and Writes during the process.
func (service *Service) localFileAvailability(filename string, response *Response) error {
	fileSysLock.RLock()
	video, ok := localFileSys.Files[filename]
	if ok {
		fmt.Println("File ", filename, " is available")
		response.Avail = true
		response.SegNums = video.SegNums
		response.SegsAvail = video.SegsAvail
	} else {
		fmt.Println("File ", filename, " is unavailable")
		response.Avail = false
	}
	fileSysLock.RUnlock()
	return nil
}

// (service *Service) sendFileSegment(filename string, segment *VidSegment) error <--
// ----------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to an rpc Call for a particular segment of a file. It first looks up and checks if the video file is available. If the
// file is available, it continues on to see if the segment is available. If the segment is available, it returns a response with the VidSegment.
// In case of unavailability, it will either return an error saying "File Unavailable." or "Segment unavailable." depending on what was unavailable.
// The method locks the local filesystem for all Reads and Writes during the process.
func (service *Service) getFileSegment(filename string, segment *VidSegment) error {
	var seg VidSegment
	fileSysLock.RLock()
	video, ok := localFileSys.Files["filename"]
	if ok {
		seg, ok = video.Segments[segment.Id]
		if ok {
			segment.Body = seg.Body
		} else {
			fileSysLock.Unlock()
			return errors.New("Segment unavailable.")
		}
	} else {
		return errors.New("File unavailable.")
	}
	fileSysLock.Unlock()
	return nil
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// OUTBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// checkFileAvailability(filename string, nodeService *rpc.Client) (bool, int64, map[int]int64)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to check if they have a video available. If it is available on the node,
// if no errors occur in the Call, the method checks the response to see if the file is available. If it is, it reads the response
// to obtain the map of segments and the total number of segments of the file
func checkFileAvailability(filename string, nodeService *rpc.Client) (bool, int64, map[int]int64) {
	var response Response
	var segNums int64
	var segsAvail map[int]int64
	err := nodeService.Call("nodeService.localFileAvailability", filename, &response)
	checkError(err)
	if response.Avail == true {
		segNums = response.SegNums
		segsAvail = response.SegsAvail
		return true, segNums, segsAvail
	} else {
		return false, 0, nil
	}
}

// getVideoSegment(filename string, segId int, nodeService *rpc.Client) (bool, int64, map[int]int64)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
func getVideoSegment(filename string, segId int, nodeService *rpc.Client) {
	var vidSeg VidSegment
	vidSeg.Id = segId
	err := nodeService.Call("nodeService.getFileSegment", filename, &vidSeg)
	checkError(err)
	// TODO
	// TODO
	// TODO
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// OTHER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// saveSegToFileSys()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
func saveSegToFileSys() {
	// TODO
	// TODO
	// TODO
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// HELPER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
}

func debug(str string) {
	fmt.Println(str)
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// MAIN METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func main() {

}
