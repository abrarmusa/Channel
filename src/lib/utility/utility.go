package utility

// package clientstream

import (
	"../colorprint"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"regexp"
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
	sync.RWMutex
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

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// HELPER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// CheckError(err error)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints error message into console in red
// -------------------
// INSTRUCTIONS:
// -------------------
// call utility.CheckError({ERROR})
func CheckError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		os.Exit(-1)
	}
}

// SaveFileInfoToJson()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method saves file information into a json file in the filesys folder
// -------------------
// INSTRUCTIONS:
// -------------------
// call utility.SaveFileInfoToJson({YOUR JSON STRUCT AS A BYTE ARRAY}, {THE DIRECTORY PATH OF THE FILESYSTEM DIRECTORY})
func SaveFileInfoToJson(jsondata []byte, dirPath string) {
	jsonFile, err := os.Create(dirPath + "/localFiles.json")
	CheckError(err)
	jsonFile.Write(jsondata)
	jsonFile.Close()
}

// ValidIP(ipAddress string, field string) bool
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Checks if the ip provided is valid. Accepts only the port as well eg. :3000 although in this case
// it assumes the localhost ip address
// -------------------
// INSTRUCTIONS:
// -------------------
// call utility.ValidIP("{YOUR IP ADDRESS STRING}", "{THE IP FORMAT STRING FOR YOUR OUTPUT")
func ValidIP(ipAddress string, field string) bool {
	re, _ := regexp.Compile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\:[0-9]{1,5}|\:[0-9]{1,5}`)
	if re.MatchString(ipAddress) {
		return true
	}
	fmt.Println("\x1b[31;1mError: "+field+":"+ipAddress, "is not in the correct format\x1b[0m")
	return false
}

// PrintFileSysTable()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints out the list of locally available files and their paths and sizes
// -------------------
// INSTRUCTIONS:
// -------------------
// call utility.PrintFileSysTable("{DIRECTORY PATH OF YOUR FILESYSTEM DIRECTORY}"")
func PrintFileSysTable(dirPath string) {
	locFiles, err := ioutil.ReadFile(dirPath + "/localFiles.json")
	CheckError(err)
	files := make([]File, 0)

	var fpaths FilePath
	fpaths.Files = files
	err = json.Unmarshal(locFiles, &fpaths)
	CheckError(err)
	colorprint.Info("LOCALLY AVAILABLE")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("   SL   |              NAME              |              DIRECTORY PATH              |       SIZE       |")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	for i, value := range fpaths.Files {
		file, err := os.Open(value.Path)
		CheckError(err)
		fi, err := file.Stat()
		CheckError(err)
		fmt.Printf("  %4d  |%30s  |%40s  | %13.2f kb |\n", (i + 1), value.Name, value.Path, float64(fi.Size())/float64(1024))
	}
	fmt.Println("--------------------------------------------------------------------------------------------------------")

}
