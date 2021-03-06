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

// This type just holds an integer to use for registering the RPC Service
type Service int

//type FTService int

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
var rpcAddress string
var progLock *sync.RWMutex
var filePaths utility.FilePath
var nodeName string
var prog Progress = Progress{
	show: true,
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// INBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

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

// func (this *FTService) SaveToDestination(folderFilePath string, fileData []byte) error {
// 	pathArr := strings.Split(folderFilePath, " ")
// 	path := "FFMPEG/NodesData/" + nodeName + "/" + pathArr[0] + "/" + pathArr[1]
// 	err := ioutil.WriteFile(path, fileData, 0644)
// 	utility.CheckError(err)
// 	return err
// }

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
		fmt.Println(segReq.SegmentId)
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

// This method answers to an rpc to save a video segment locally into the local filesystem
func (service *Service) ReceiveFileSegment(seqStruct *utility.SeqStruct, segment *utility.VidSegment) error {
	fmt.Println("SS")
	filename := seqStruct.Filename
	t := time.Now().String()
	outputstr := ""
	colorprint.Debug("------------------------------------------------------------------")
	colorprint.Debug(">> " + t + "  <<")
	colorprint.Debug("INBOUND RPC REQUEST: Receiving video segment for " + seqStruct.Filename)
	// localFileSys.Lock()
	filemgmt.AddVidSegIntoFileSys(filename, int64(seqStruct.SegNums), seqStruct.Segment, &localFileSys)
	outputstr += ("\nSegment " + strconv.Itoa(seqStruct.Segment.Id) + " received for " + filename)
	colorprint.Warning(outputstr)
	//localFileSys.Unlock()
	colorprint.Warning("Video Segment saved on the node")
	return nil
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// OUTBOUND RPC CALL METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

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
func GetVideoSegment(fname string, segNums int64, segId int, nodeAdd string) utility.VidSegment {
	nodeService, err := rpc.Dial(consts.TransProtocol, nodeAdd)
	// utility.CheckError(err)
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

	prog.Lock()
	prog.show = true
	prog.Unlock()
	segReq := &utility.ReqStruct{
		Filename:  fname,
		SegmentId: segId,
	}
	var vidSeg utility.VidSegment
	vidSeg.Id = segId
	err = nodeService.Call("Service.GetFileSegment", segReq, &vidSeg)
	utility.CheckError(err)
	err = nodeService.Close()
	utility.CheckError(err)
	filemgmt.AddVidSegIntoFileSys(fname, segNums, vidSeg, &localFileSys)
	return vidSeg
}

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
	fmt.Printf("\rSending segment " + strconv.Itoa(segment.Id))
	waitstr := "."
	counter, incrementer := 0, 1
	nodeService, err := rpc.Dial(consts.TransProtocol, nodeAdd)
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
	segReq := utility.SeqStruct{
		Filename:  fname,
		SegNums:   segNums,
		SegmentId: segment.Id,
		Segment:   segment,
	}
	err = nodeService.Call("Service.ReceiveFileSegment", segReq, &segment)
	utility.CheckError(err)
	err = nodeService.Close()
	utility.CheckError(err)

}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// CONNECTION METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// This method sets up the RPC connection using UDP
func setUpRPC(nodeRPC string) {
	rpcServ := new(Service)
	rpc.Register(rpcServ)
	rpcAddr, err := net.ResolveTCPAddr("tcp", nodeRPC)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	l, e := net.ListenTCP(consts.TransProtocol, rpcAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for i := 0; i >= 0; i++ {
		conn, _ := l.AcceptTCP()
		colorprint.Alert("=========================================================================================")
		colorprint.Debug("REQ " + strconv.Itoa(i) + ": ESTABLISHING RPC REQUEST CONNECTION WITH " + conn.LocalAddr().String())
		go rpc.ServeConn(conn)
		colorprint.Blue("REQ " + strconv.Itoa(i) + ": Request Served")
		colorprint.Alert("=========================================================================================")
		defer conn.Close()
	}
	l.Close()

	// rpcServ := new(FTService)
	// rpc.Register(rpcServ)
	// rpcAddr, err := net.ResolveTCPAddr("tcp", nodeRPC)
	// if err != nil {
	// 	log.Fatal("listen error:", err)
	// }
	// l, e := net.ListenTCP(consts.TransProtocol, rpcAddr)
	// if e != nil {
	// 	log.Fatal("listen error:", e)
	// }
	// for i := 0; i >= 0; i++ {
	// 	conn, _ := l.AcceptTCP()
	// 	colorprint.Alert("=========================================================================================")
	// 	colorprint.Debug("REQ " + strconv.Itoa(i) + ": ESTABLISHING RPC REQUEST CONNECTION WITH " + conn.LocalAddr().String())
	// 	rpc.ServeConn(conn)
	// 	colorprint.Blue("REQ " + strconv.Itoa(i) + ": Request Served")
	// 	colorprint.Alert("=========================================================================================")
	// 	//defer conn.Close()
	// }
	// l.Close()

}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// IO METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// This method responds to the user input requests
func Instr() {
	colorprint.Blue("Listening on " + rpcAddress + " for incoming RPC calls")
	var input string
	var fname string
	for i := 0; i >= 0; i++ {
		fmt.Println()
		colorprint.Info(">>>> Please enter a command")
		fmt.Scan(&input)
		cmd := input
		if input == "get" {
			getHelper(rpcAddress, input, fname, cmd)
		} else if input == "list" {

		} else if input == "help" {
			ui.Help()
		} else if input == "play" {
			playHelper(fname)
		}
	}
}

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
			// TODO
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

// This method starts up the transfer rpc and also initializes the filesystem
func Initialize(nodeRPC string, name string) *utility.FileSys {
	if !utility.ValidIP(nodeRPC, "[node RPC ip:port]") {
		colorprint.Alert("Please provide a valid IP string.")
		return nil
	}
	// ========================================
	progLock = &sync.RWMutex{}
	// ========================================
	rpcAddress = nodeRPC
	go setUpRPC(nodeRPC)
	nodeName = name
	localFileSys = utility.FileSys{
		Id:    1,
		Files: make(map[string]utility.Video),
	}
	return &localFileSys
}
