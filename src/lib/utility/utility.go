package utility

// package clientstream

import (
	"../../consts"
	"fmt"
	"github.com/fatih/color"
	"os"
	"regexp"
	"sync"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//  STRUCTS & TYPES
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


// This struct holds a particular part of a video file. Id refers to the segment id and the body refers to the bytecount of the actual video bytes
type VidSegment struct {
	Id   int
	Body []byte
}


// This struct holds the info for obtaining a video segment
type ReqStruct struct {
	Filename  string
	SegmentId int
}


// This struct holds the fields for sending a video segment to be saved onto a node
type SeqStruct struct {
	Filename  string
	SegNums   int
	SegmentId int
	Segment   VidSegment
}

--
// This struct holds the VidSegments of a particular video file. SegNums refers to the total number of segments in the entire video.
// This colorprint.Info is used to reorder the segments when playing the video stream.
type Video struct {
	Name      string
	SegNums   int64
	SegsAvail []int64
	Segments  map[int]VidSegment
}


// This struct represents the local FileSystem to hold the Video's. Each FileSys object has an id and a map of Files with the keys to the map being the
// actual filename.
type FileSys struct {
	Id    int
	Files map[string]Video
	sync.RWMutex
}


// This struct represents the Files as their names and the directory that they are located in
type File struct {
	Name      string  `json:"name"`
	Path      string  `json:"dir"`
	SegNums   int64   `json:"segnums"`
	SegsAvail []int64 `json:"segsavail"`
}


// This struct represents the filenames and their directory paths from which we are to read the files to be processed
type FilePath struct {
	Files []File `json:"Files"`
}


// This struct represents the response that an RPC call will write to. It is used to check if a node has a particular file and if it does, which parts
// of that file it has in its local filesystem.
type Response struct {
	Avail     bool
	SegNums   int64
	SegsAvail []int64
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// HELPER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


// Prints error message into console in red
func CheckError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		os.Exit(-1)
	}
}

// This method saves file information into a json file in the filesys folder
func SaveFileInfoToJson(jsondata []byte) {
	jsonFile, err := os.Create(consts.DirPath + "/localFiles.json")
	CheckError(err)
	jsonFile.Write(jsondata)
	jsonFile.Close()
}

// Checks if the ip provided is valid. Accepts only the port as well eg. :3000 although in this case
// it assumes the localhost ip address
func ValidIP(ipAddress string, field string) bool {
	re, _ := regexp.Compile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\:[0-9]{1,5}|\:[0-9]{1,5}`)
	if re.MatchString(ipAddress) {
		return true
	}
	fmt.Println("\x1b[31;1mError: "+field+":"+ipAddress, "is not in the correct format\x1b[0m")
	return false
}
