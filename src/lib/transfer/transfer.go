package transfer

import (
	"../../consts"
	"../colorprint"
	"../filemgmt"
	"../player"
	"../ui"
	"../utility"
	"errors"
	"fmt"
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

// (service *Service) ReceiveFileSegment(seqStruct utility.SeqStruct, segment *utility.VidSegment)
// ----------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method answers to an rpc to save a video segment locally into the local filesystem
func (service *Service) ReceiveFileSegment(seqStruct *utility.SeqStruct, segment *utility.VidSegment) string {
	filename := seqStruct.Filename
	t := time.Now().String()
	outputstr := ""
	colorprint.Debug("------------------------------------------------------------------")
	colorprint.Debug(">> " + t + "  <<")
	colorprint.Debug("INBOUND RPC REQUEST: Receiving video segment for " + seqStruct.Filename)
	filemgmt.AddVidSegIntoFileSys(filename, int64(seqStruct.SegNums), *segment, &localFileSys)
	outputstr += ("\nSegment " + strconv.Itoa(segment.Id) + " received for " + filename)
	colorprint.Warning(outputstr)
	localFileSys.Unlock()
	return "Video Segment saved on the node"
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
// Call transfer.CheckFileAvailability("{FILENAME STRING}, "{THE ADDRESS OF THE NODE FOR THE RPC CALL}", {*rpc.Client}")
//
// E.g code:
//
// avail, segNums, segsAvail := transfer.CheckFileAvailability("sample.mp4", ":3000")
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
// Call transfer.GetVideoSegment("sample.mp4", 45, ":3000")
//
// E.g code:
//
// segNums := 100
// vidMap := make(map[int]utility.VidSegment)
// for i := 0; i < segNums; i++ {
// 		vidMap[i] = transfer.GetVideoSegment("sample.mp4", 45, ":3000")
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

// SendVideoSegment(fname string, nodeAddr string, segment utility.VidSegment)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method sends a utility.VidSegment to another node for saving
// -------------------
// INSTRUCTIONS:
// -------------------
// Call transfer.GetVideoSegment("sample.mp4", 45, ":3000")
//
// E.g code:
//
// var segment utility.VidSegment
// vidMap := make(map[int]utility.VidSegment)
// for i := 0; i < segNums; i++ {
// 		vidMap[i] = transfer.SendVideoSegment("sample.mp4", ":3000", segment)
// }
//
func SendVideoSegment(fname string, nodeAdd string, segNums int, segment utility.VidSegment) {
	nodeService, err := rpc.Dial(consts.TransProtocol, nodeAdd)
	utility.CheckError(err)
	segReq := utility.SeqStruct{
		Filename:  fname,
		SegNums:   segNums,
		SegmentId: segment.Id,
	}
	err = nodeService.Call("Service.ReceiveFileSegment", segReq, &segment)
	utility.CheckError(err)
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
	avail, _, _ := CheckFileAvailability(fname, nodeAddr)
	if avail && (cmd == "get") {
		colorprint.Info(">>>> Would you like to get the file from the node[" + nodeRPC + "]?(y/n)")
		fmt.Scan(&input)
		colorprint.Debug("<<<< " + input)
		if input == "y" {

			// saveSegsToFileSys(nodeAddr, segNums, fname)
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
// This method starts up the transfer rpc and also initializes the filesystem
func Initialize(nodeRPC string) *utility.FileSys {
	if !utility.ValidIP(nodeRPC, "[node RPC ip:port]") {
		colorprint.Alert("Please provide a valid IP string.")
		return nil
	}
	// ========================================
	// localFileSys = &sync.RWMutex{}
	progLock = &sync.RWMutex{}
	// ========================================
	go setUpRPC(nodeRPC)
	localFileSys = utility.FileSys{
		Id:    1,
		Files: make(map[string]utility.Video),
	}
	return &localFileSys
}
