package fileshare

// package clientstream

import (
	"../../consts"
	"../colorprint"
	"../player"
	"../ui"
	"../utility"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//  STRUCTS & TYPES
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// --> Service <---
// ----------------
// DESCRIPTION:
// -------------------
// This type just holds an integer to use for registering the RPC Service
type Service int

// --> VidSegment <---
// -------------------
// DESCRIPTION:
// -------------------
// This struct holds a bool and a lock to determine when a rpc.Dial method is waiting for ports to clear up
type Progress struct {
	show bool
	sync.RWMutex
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// GLOBAL VARS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
var localFileSys utility.FileSys
var progLock *sync.RWMutex
var filePaths utility.FilePath
var prog Progress = Progress{
	show: true,
}

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
// -------------------
// INSTRUCTIONS:
// -------------------
// THIS FUNCTION IS NEVER CALLED BY ANY MAIN OPERATION. THIS IS THE RPC FUNCTION TO RESPOND FOR FILE AVAILABILITY.
// The returned response is either an error or a *utility.Response if filled out. If the response is available,
// response.Avail == true and all its other fields will be filled in. If the response is unavailable, it will be false
// and an error will be returned
//
// E.g code:
//
// var response utility.Response
// var segNums int64
// var segsAvail []int64
// err := nodeService.Call("Service.LocalFileAvailability", filename, &response)
// utility.CheckError(err)
//
func (service *Service) LocalFileAvailability(filename string, response *utility.Response) error {
	colorprint.Debug("INBOUND RPC REQUEST: Checking utility.File Availability " + filename)
	localFileSys.RLock()
	video, ok := localFileSys.Files[filename]
	if ok {
		colorprint.Info("utility.File " + filename + " is available")
		var segstr string = ""
		for _, value := range video.SegsAvail {
			segstr += strconv.FormatInt(value, 10) + " "
		}
		colorprint.Info("Locally available segments are " + segstr + "out of " + strconv.FormatInt(video.SegNums, 10))
		colorprint.Debug("INBOUND RPC REQUEST COMPLETED")
		response.Avail = true
		response.SegNums = video.SegNums
		response.SegsAvail = video.SegsAvail
	} else {
		colorprint.Alert("utility.File " + filename + " is unavailable")
		response.Avail = false
		return errors.New("utility.File " + filename + " is unavailable")

	}
	localFileSys.RUnlock()
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
	t := time.Now().String()
	colorprint.Debug("------------------------------------------------------------------")
	colorprint.Debug(">> " + t + "  <<")
	colorprint.Debug("INBOUND RPC REQUEST: Sending video segment for " + segReq.Filename)
	var seg utility.VidSegment
	localFileSys.RLock()
	outputstr := ""
	video, ok := localFileSys.Files[segReq.Filename]
	if ok {
		outputstr += ("Node is asking for segment no. " + strconv.Itoa(segReq.SegmentId) + " for " + segReq.Filename)
		_, ok := (video.Segments[1])
		seg, ok = video.Segments[segReq.SegmentId]
		if ok {
			segment.Body = seg.Body
		} else {
			outputstr += ("\nSegment " + strconv.Itoa(segReq.SegmentId) + " unavailable for " + segReq.Filename)
			return errors.New("Segment unavailable.")
			localFileSys.Unlock()
		}
	} else {

		return errors.New("utility.File unavailable.")
		localFileSys.Unlock()
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
// -------------------
// INSTRUCTIONS:
// -------------------
// Call fileshare.CheckFileAvailability("{FILENAME STRING}, "{THE ADDRESS OF THE NODE FOR THE RPC CALL}", {*rpc.Client}")
//
// E.g code:
//
// avail, segNums, segsAvail := fileshare.CheckFileAvailability("sample.mp4", ":3000")
//
func CheckFileAvailability(filename string, nodeadd string) (bool, int64, []int64) {
	colorprint.Debug("OUTBOUND REQUEST: Check utility.File Availability")
	var response utility.Response
	var segNums int64
	var segsAvail []int64
	nodeService, err := rpc.Dial(consts.TransProtocol, nodeadd)
	utility.CheckError(err)
	err = nodeService.Call("Service.LocalFileAvailability", filename, &response)
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

// GetVideoSegment(fname string, segId int, nodeAdd string) utility.VidSegment
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
// -------------------
// INSTRUCTIONS:
// -------------------
// Call fileshare.GetVideoSegment("sample.mp4", 45, ":3000")
//
// E.g code:
//
// segNums := 100
// vidMap := make(map[int]utility.VidSegment)
// for i := 0; i < segNums; i++ {
// 		vidMap[i] = fileshare.GetVideoSegment("sample.mp4", 45, ":3000")
// }
//
func GetVideoSegment(fname string, segId int, nodeAdd string) utility.VidSegment {
	nodeService, err := rpc.Dial(consts.TransProtocol, nodeAdd)
	waitstr := "."
	counter, incrementer := 0, 1
	if err != nil {
		prog.Lock()
		prog.show = false
		prog.Unlock()
	}
	for err != nil {
		counter++
		if (counter % 100) == 0 {
			fmt.Printf("\rAll ports are blocked. Waiting for port to clear%s                          ", waitstr)
			waitstr += "."
			incrementer++
			if (incrementer % 300) == 0 {
				waitstr += "."
			}
		}
		nodeService, err = rpc.Dial(consts.TransProtocol, nodeAdd)
	}
	// utility.CheckError(err)
	prog.Lock()
	prog.show = true
	prog.Unlock()
	segReq := &utility.ReqStruct{
		Filename:  fname,
		SegmentId: segId,
	}
	var vidSeg utility.VidSegment
	err = nodeService.Call("Service.GetFileSegment", segReq, &vidSeg)
	utility.CheckError(err)
	err = nodeService.Close()
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
func saveSegsToFileSys(nodeAddr string, segNums int64, fname string) {
	var newVid utility.Video
	vidMap := make(map[int]utility.VidSegment)
	var segsAvail []int64
	progstr := "="
	counter2, counter3, altc, downloadstr := 1, 1, (segNums / consts.Factor), 0
	quit := make(chan int)
	go func() {
		c := time.Tick(1 * time.Second)
		for _ = range c {
			select {
			default:
				prog.RLock()
				if prog.show == true {
					downloadstr = ui.ProgressBar(progstr, counter3, downloadstr, int(segNums))
				}
				prog.RUnlock()
			case <-quit:
				return
			}
		}

	}()
	for i := 1; i <= int(segNums); i++ {
		vidMap[i] = GetVideoSegment(fname, i, nodeAddr)
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
	newVid = utility.Video{
		Name:      fname,
		SegNums:   segNums,
		SegsAvail: segsAvail,
		Segments:  vidMap,
	}
	localFileSys.Lock()
	localFileSys.Files[fname] = newVid
	localFileSys.Unlock()
	fmt.Println()

	colorprint.Warning("Saving file info into filesystem table")
	localFileSys.RLock()
	vid := localFileSys.Files[fname]
	localFileSys.RUnlock()
	pathname := writeToFileHelper(fname, vid)
	newFile := utility.File{
		Name: fname,
		Path: pathname,
	}
	filePaths.Files = append(filePaths.Files, newFile)
	jsondata, err := json.Marshal(filePaths)
	utility.CheckError(err)
	utility.SaveFileInfoToJson(jsondata, consts.DirPath)
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
	str := consts.DirPath + "/downloaded/" + fname
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
	locFiles, err := ioutil.ReadFile("../filesys/localFiles.json")
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
	fmt.Println("========================    PROCESSING LOCAL FILES FOR SHARING    ========================")
	fmt.Println("==========================================================================================")
	for index, value := range filePaths.Files {

		dat, err := ioutil.ReadFile(value.Path)
		utility.CheckError(err)
		colorprint.Info("---------------------------------------------------------------------------")
		colorprint.Info(strconv.Itoa(index+1) + ": Processing " + value.Name + " at " + value.Path + " with " + strconv.Itoa(len(dat)/consts.Bytecount) + " segments.")
		segsAvail, vidMap := convByteArrayToSeg(dat)

		vid := utility.Video{
			Name:      value.Name,
			SegNums:   int64(len(dat) / consts.Bytecount),
			SegsAvail: segsAvail,
			Segments:  vidMap,
		}
		localFileSys.Lock()
		localFileSys.Files[value.Name] = vid
		localFileSys.Unlock()
		colorprint.Info("Completed Processing " + value.Name + " at " + value.Path)
		colorprint.Info("---------------------------------------------------------------------------")

	}
	fmt.Println("===============================    PROCESSING COMPLETE    ================================\n\n\n")

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
	altc := (blen / int(consts.Factor))
	for index, element := range bytes {
		eightBSeg = append(eightBSeg, element)
		if counter == consts.Bytecount {
			counter = 0
			vidSeg = utility.VidSegment{
				Id:   ((index % consts.Bytecount) + 1),
				Body: eightBSeg,
			}
			vidmap[((index / consts.Bytecount) + 1)] = vidSeg
			segsAvail = append(segsAvail, int64(((index / consts.Bytecount) + 1)))
			eightBSeg = []byte{}
		}
		counter++
		counter2++
		counter3++
		if counter2 == altc {
			progstr += "~"
			fmt.Printf("\r|%s|  - %d%%", progstr, ((counter3*100)/blen + 1))
			counter2 = 0
		}
	}
	fmt.Println()
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
	l, e := net.Listen(consts.TransProtocol, nodeRPC)
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
		defer conn.Close()
	}

}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// IO METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Instr(nodeRPC string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to the user input requests
func Instr(nodeRPC string) {
	var input string
	var fname string
	for i := 0; i >= 0; i++ {
		fmt.Println()
		colorprint.Info(">>>> Please enter a command")
		fmt.Scan(&input)
		cmd := input
		if input == "get" {
			getHelper(nodeRPC, input, fname, cmd)
		} else if input == "list" {
			utility.PrintFileSysTable(consts.DirPath)
		} else if input == "help" {
			ui.Help()
		} else if input == "play" {
			playHelper(fname)
		}
	}
}

// getHelper(nodeRPC string, nodeUDP string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method responds to the user input request for "get"
func getHelper(nodeRPC string, input string, fname string, cmd string) {

	colorprint.Info(">>>> Please enter the name of the file that you would like to obtain")
	fmt.Scan(&fname)
	colorprint.Debug("<<<< " + fname)
	colorprint.Info(">>>> Please enter the address of the node you want to connect to")
	fmt.Scan(&input)
	colorprint.Debug("<<<< " + input)
	nodeAddr := input
	// Connect to utility.Service via RPC // returns *Client, err
	avail, segNums, _ := CheckFileAvailability(fname, nodeAddr)
	if avail && (cmd == "get") {
		colorprint.Info(">>>> Would you like to get the file from the node[" + nodeRPC + "]?(y/n)")
		fmt.Scan(&input)
		colorprint.Debug("<<<< " + input)
		if input == "y" {

			saveSegsToFileSys(nodeAddr, segNums, fname)
		}
	}
}

// playHelper(fname string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method plays the requested video file. Visit localhost:8080
func playHelper(fname string) {

	colorprint.Info(">>>> Please enter the name of the file that you would like to play")
	fmt.Scan(&fname)
	colorprint.Debug("<<<< " + fname)
	player.Filename = fname
	for _, element := range filePaths.Files {
		if element.Name == fname {
			player.Filepath = element.Path
		}
	}

	player.Run()
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// MAIN METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// playHelper(fname string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method starts up the filesystem
func FileSysStart(nodeRPC string) {
	if !utility.ValidIP(nodeRPC, "[node RPC ip:port]") {
		colorprint.Alert("Please provide a valid IP string.")
		return
	}
	// ========================================
	// localFileSys = &sync.RWMutex{}
	progLock = &sync.RWMutex{}
	processLocalVideosIntoFileSys()
	// ========================================
	go setUpRPC(nodeRPC)
}
