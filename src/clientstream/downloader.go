package main

// package clientstream

import (
	"./colorprint"
	"./utility"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

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
var localFileSys utility.FileSys
var fileSysLock *sync.RWMutex
var bytecount int = 1024
var filePaths utility.FilePath

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// INBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// (service *Service) localFileAvailability(filename string, response *utilty.Response) error
// -----------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to an rpc Call for a particular segment of a file. It first looks up and checks if the video file is available. If
// the file is available, it continues on to see if the segment is available. If the segment is available, it returns a response with the utility.VidSegment.
// In case of unavailability, it will either return an error saying "utility.File Unavailable." or "Segment unavailable." depending on what was unavailable.
// The method locks the local filesystem for all Reads and Writes during the process.
func (service *Service) LocalFileAvailability(filename string, response *utility.Response) error {
	colorprint.Debug("INBOUND RPC REQUEST: Checking utility.File Availability " + filename)
	fileSysLock.RLock()
	video, ok := localFileSys.Files[filename]
	if ok {
		colorprint.Info("utility.File " + filename + " is available")
		colorprint.Debug("INBOUND RPC REQUEST COMPLETED")
		response.Avail = true
		response.SegNums = video.SegNums
		response.SegsAvail = video.SegsAvail
	} else {
		colorprint.Alert("utility.File " + filename + " is unavailable")
		response.Avail = false
	}
	fileSysLock.RUnlock()
	return nil
}

// (service *Service) GetFileSegment(segReq *utility.ReqStruct, segment *utility.VidSegment) error <--
// ----------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to an rpc Call for a particular segment of a file. It first looks up and checks if the video file is available. If the
// file is available, it continues on to see if the segment is available. If the segment is available, it returns a response with the utility.VidSegment.
// In case of unavailability, it will either return an error saying "utility.File Unavailable." or "Segment unavailable." depending on what was unavailable.
// The method locks the local filesystem for all Reads and Writes during the process.
func (service *Service) GetFileSegment(segReq *utility.ReqStruct, segment *utility.VidSegment) error {
	colorprint.Debug("INBOUND RPC REQUEST: Sending video segment for " + segReq.Filename)
	var seg utility.VidSegment
	fileSysLock.RLock()
	outputstr := ""
	video, ok := localFileSys.Files[segReq.Filename]
	if ok {
		outputstr += ("\nNode is asking for segment no. " + strconv.Itoa(segReq.SegmentId) + " for " + segReq.Filename)
		_, ok := (video.Segments[1])
		if ok {
			outputstr += ("\nSeg 1 available")
		}
		seg, ok = video.Segments[segReq.SegmentId]
		if ok {
			segment.Body = seg.Body
		} else {
			outputstr += ("\nSegment " + strconv.Itoa(segReq.SegmentId) + " unavailable for " + segReq.Filename)
			return errors.New("Segment unavailable.")
			fileSysLock.Unlock()
		}
	} else {

		return errors.New("utility.File unavailable.")
		fileSysLock.Unlock()
	}
	colorprint.Warning(outputstr)
	return nil
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// OUTBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// checkFileAvailability(filename string, nodeadd string, nodeService *rpc.Client) (bool, int64, []int64)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to check if they have a video available. If it is available on the node,
// if no errors occur in the Call, the method checks the response to see if the file is available. If it is, it reads the response
// to obtain the map of segments and the total number of segments of the file
func checkFileAvailability(filename string, nodeadd string, nodeService *rpc.Client) (bool, int64, []int64) {
	colorprint.Debug("OUTBOUND REQUEST: Check utility.File Availability")
	var response utility.Response
	var segNums int64
	var segsAvail []int64
	err := nodeService.Call("Service.LocalFileAvailability", filename, &response)
	utility.CheckError(err)
	colorprint.Debug("OUTBOUND REQUEST COMPLETED")
	if response.Avail == true {
		fmt.Println("utility.File:", filename, " is available")
		segNums = response.SegNums
		segsAvail = response.SegsAvail
		return true, segNums, segsAvail
	} else {
		fmt.Println("utility.File:", filename, " is not available on node["+""+"].")
		return false, 0, nil
	}
}

// getVideoSegment(fname string, segId int, nodeService *rpc.Client) utility.VidSegment
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
func getVideoSegment(fname string, segId int, nodeService *rpc.Client) utility.VidSegment {
	segReq := &utility.ReqStruct{
		Filename:  fname,
		SegmentId: segId,
	}
	var vidSeg utility.VidSegment
	err := nodeService.Call("Service.GetFileSegment", segReq, &vidSeg)
	utility.CheckError(err)
	return vidSeg
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// OTHER METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// saveSegsToFileSys(service *rpc.Client, segNums int64, fname string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
func saveSegsToFileSys(service *rpc.Client, segNums int64, fname string) {
	var newVid utility.Video
	vidMap := make(map[int]utility.VidSegment)
	var segsAvail []int64
	progstr := "="
	counter2, counter3, altc, downloadstr := 1, 1, (segNums / 100), 0
	quit := make(chan int)
	go func() {
		for i := 1; i <= int(segNums); i++ {
			vidMap[i] = getVideoSegment(fname, i, service)
			segsAvail = append(segsAvail, int64(i))
			counter2++
			counter3++
			downloadstr++
			if counter2 == int(altc) {
				progstr += "="
				counter2 = 0
			}
			if i == int(segNums) {
				quit <- 1
			}
		}
	}()
forTimer:
	for range time.Tick(1 * time.Second) {
		progress := ((counter3 * 100) / int(segNums))
		fmt.Printf("\rDownload Speed: %.1f MB/s [%s]  - %d%%", float64(downloadstr)/float64(bytecount), progstr, progress)
		downloadstr = 0
		select {
		case <-quit:
			break forTimer
		default:
			continue
		}
	}
	newVid = utility.Video{
		Name:      fname,
		SegNums:   segNums,
		SegsAvail: segsAvail,
		Segments:  vidMap,
	}
	fileSysLock.Lock()
	localFileSys.Files[fname] = newVid
	fileSysLock.Unlock()
	fmt.Println()

	colorprint.Warning("Saving file info into filesystem table")
	fileSysLock.RLock()
	vid := localFileSys.Files[fname]
	fileSysLock.RUnlock()
	pathname := writeToFileHelper(fname, vid)
	newFile := utility.File{
		Name: fname,
		Path: pathname,
	}
	filePaths.Files = append(filePaths.Files, newFile)
	jsondata, err := json.Marshal(filePaths)
	utility.CheckError(err)
	utility.SaveFileInfoToJson(jsondata)
	colorprint.Info(fname + " saved into file system. File is located at " + pathname + ".")

}

// writeToFileHelper()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method writes the downloaded file into a file of its own in the filesys/downloaded folder
func writeToFileHelper(fname string, video utility.Video) string {
	var data []byte
	for i := 1; i < int(video.SegNums); i++ {
		for j := 0; j < len(video.Segments[i].Body); j++ {
			data = append(data, video.Segments[i].Body[j])
		}
	}
	str := "./filesys/downloaded/"
	str += fname
	err := ioutil.WriteFile(str, data, 0777)
	utility.CheckError(err)
	return str

}

// processLocalVideosIntoFileSys()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method loads up a local json file to see which files are available in the local file system. Once
// the read has been completed, the files are then processed into the utility.utility.FileSys map accordingly
func processLocalVideosIntoFileSys() {
	locFiles, err := ioutil.ReadFile("./filesys/localFiles.json")
	utility.CheckError(err)
	files := make([]utility.File, 0)

	filePaths.Files = files
	err = json.Unmarshal(locFiles, &filePaths)
	utility.CheckError(err)
	// Initialize local file system
	localFileSys = utility.FileSys{
		Id:    1,
		Files: make(map[string]utility.Video),
	}
	for _, value := range filePaths.Files {
		colorprint.Info("Processing " + value.Name + " at " + value.Path)
		dat, err := ioutil.ReadFile(value.Path)
		utility.CheckError(err)
		colorprint.Info("---------------------------------------------------------------------------")
		colorprint.Info("utility.Video:" + value.Name + " has " + strconv.Itoa(len(dat)/bytecount) + " segments.")
		segsAvail, vidMap := convByteArrayToSeg(dat)

		vid := utility.Video{
			Name:      value.Name,
			SegNums:   int64(len(dat) / bytecount),
			SegsAvail: segsAvail,
			Segments:  vidMap,
		}
		utility.PrintAvSegs(segsAvail)
		fileSysLock.Lock()
		localFileSys.Files[value.Name] = vid
		fileSysLock.Unlock()
		colorprint.Info("Completed Processing " + value.Name + " at " + value.Path)
		colorprint.Info("---------------------------------------------------------------------------")
	}

}

// convByteArrayToSeg(bytes []byte) ([]int64, map[int]utility.VidSegment)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Converts the byte array from a video files into utility.Video Segments.
func convByteArrayToSeg(bytes []byte) ([]int64, map[int]utility.VidSegment) {
	vidmap := make(map[int]utility.VidSegment)
	var segsAvail []int64
	var vidSeg utility.VidSegment
	var eightBSeg []byte
	counter, counter2, counter3 := 1, 1, 1
	progstr := "="
	blen := len(bytes)
	altc := (blen / 100)
	for index, element := range bytes {
		eightBSeg = append(eightBSeg, element)
		if counter == bytecount {
			counter = 0
			vidSeg = utility.VidSegment{
				Id:   ((index % bytecount) + 1),
				Body: eightBSeg,
			}
			vidmap[((index / bytecount) + 1)] = vidSeg
			segsAvail = append(segsAvail, int64(((index / bytecount) + 1)))
			eightBSeg = []byte{}
		}
		counter++
		counter2++
		counter3++
		if counter2 == altc {
			progstr += "="
			fmt.Printf("\r[%s]  - %d%%", progstr, ((counter3*100)/blen + 1))
			counter2 = 0
		}
	}
	fmt.Println()
	colorprint.Debug("SEGMENTS PROCESSED: " + strconv.Itoa((len(segsAvail))))
	return segsAvail, vidmap
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// CONNECTION METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// setUpRPC(nodeRPC string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method sets up the RPC connection using UDP
func setUpRPC(nodeRPC string) {
	rpcServ := new(Service)
	rpc.Register(rpcServ)
	l, e := net.Listen("tcp", nodeRPC)

	if e != nil {
		log.Fatal("listen error:", e)
	}
	colorprint.Blue("Listening on " + nodeRPC + " for incoming RPC calls")
	for i := 0; i >= 0; i++ {
		conn, _ := l.Accept()
		colorprint.Alert("=========================================================================================")
		colorprint.Debug("REQ " + strconv.Itoa(i) + ": ESTABLISHING RPC REQUEST CONNECTION WITH " + conn.LocalAddr().String())
		go rpc.ServeConn(conn)
		colorprint.Blue("REQ " + strconv.Itoa(i) + ": Request Served")
		colorprint.Alert("=========================================================================================")
	}

}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// IO METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// instr(nodeRPC string, nodeUDP string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to the user input requests
func instr(nodeRPC string, nodeUDP string) {
	var input string
	var fname string
	for i := 0; i >= 0; i++ {
		colorprint.Info(">>>> Please type in the command")
		fmt.Scan(&input)
		cmd := input
		if input == "get" {
			getHelper(nodeRPC, nodeUDP, input, fname, cmd)
		} else if input == "list" {
			utility.PrintFileSysTable()
		}
	}
}

// getHelper(nodeRPC string, nodeUDP string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to the user input request for "get"
func getHelper(nodeRPC string, nodeUDP string, input string, fname string, cmd string) {

	colorprint.Info(">>>> Please enter the name of the file that you would like to obtain")
	fmt.Scan(&fname)
	colorprint.Debug("<<<< " + fname)
	colorprint.Info(">>>> Please enter the address of the node you want to connect to")
	fmt.Scan(&input)
	colorprint.Debug("<<<< " + input)
	nodeAddr := input
	service, err := rpc.Dial("tcp", nodeAddr) // Connect to utility.Service via RPC // returns *Client, err
	utility.CheckError(err)
	avail, segNums, _ := checkFileAvailability(fname, nodeAddr, service)
	if avail && (cmd == "get") {
		colorprint.Info(">>>> Would you like to get the file from the node[" + nodeRPC + "]?(y/n)")
		fmt.Scan(&input)
		colorprint.Debug("<<<< " + input)
		if input == "y" {
			saveSegsToFileSys(service, segNums, fname)
		}
	}
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// MAIN METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func main() {
	// ========================================
	fileSysLock = &sync.RWMutex{}
	processLocalVideosIntoFileSys()
	// ========================================
	if len(os.Args) == 3 {
		nodeRPC := os.Args[1]
		nodeUDP := os.Args[2]
		if !utility.ValidIP(nodeRPC, "[node RPC ip:port]") || !utility.ValidIP(nodeUDP, "[node UDP ip:port]") {
			os.Exit(-1)
		}
		go setUpRPC(nodeRPC)
		instr(nodeRPC, nodeUDP)

	}
}
