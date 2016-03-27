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
func checkFileAvailability(filename string, nodeadd string, nodeService *rpc.Client) (bool, int64, []int64) {
	colorprint.Debug("OUTBOUND REQUEST: Check File Availability")
	var response Response
	var segNums int64
	var segsAvail []int64
	err := nodeService.Call("Service.LocalFileAvailability", filename, &response)
	checkError(err)
	colorprint.Debug("OUTBOUND REQUEST COMPLETED")
	if response.Avail == true {
		fmt.Println("File:", filename, " is available")
		segNums = response.SegNums
		segsAvail = response.SegsAvail
		return true, segNums, segsAvail
	} else {
		fmt.Println("File:", filename, " is not available on node["+""+"].")
		return false, 0, nil
	}
}

// getVideoSegment(filename string, segId int, nodeService *rpc.Client) (bool, int64, map[int]int64)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method calls an RPC method to another node to obtain a particular segment of a video
func getVideoSegment(fname string, segId int, nodeService *rpc.Client) VidSegment {
	segReq := &ReqStruct{
		Filename:  fname,
		SegmentId: segId,
	}
	var vidSeg VidSegment
	err := nodeService.Call("Service.GetFileSegment", segReq, &vidSeg)
	checkError(err)
	return vidSeg
}
