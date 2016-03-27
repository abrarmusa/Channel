package utility

// package clientstream

import (
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
// This struct holds a particular part of a video file. Id refers to the segment id and the body refers to the bytecount of the actual video bytes
type VidSegment struct {
	Id   int
	Body []byte
}

type ReqStruct struct {
	Filename  string
	SegmentId int
}

// --> Video <---
// --------------
// DESCRIPTION:
// -------------------
// This struct holds the VidSegments of a particular video file. SegNums refers to the total number of segments in the entire video.
// This colorprint.Info is used to reorder the segments when playing the video stream.
type Video struct {
	Name      string
	SegNums   int64
	SegsAvail []int64
	Segments  map[int]VidSegment
}

// --> FileSys <---
// ----------------
// DESCRIPTION
// -------------------
// This struct represents the local FileSystem to hold the Video's. Each FileSys object has an id and a map of Files with the keys to the map being the
// actual filename.
// This colorprint.Info is used to check if a node actually has the file
type FileSys struct {
	Id    int
	Files map[string]Video
}

// --> File <---
// ----------------
// DESCRIPTION
// -------------------
// This struct represents the Files as their names and the directory that they are located in
type File struct {
	Name string `json:"name"`
	Path string `json:"dir"`
}

// --> FilePath <---
// ----------------
// DESCRIPTION
// -------------------
// This struct represents the filenames and their directory paths from which we are to read the files to be processed
type FilePath struct {
	Files []File `json:"Files"`
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
	SegsAvail []int64
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
var bytecount int = 2048
