package inboundRPC

// package clientstream

import (
	"../colorprint"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"regexp"
	"strconv"
	"sync"
)

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
func (service *Service) LocalFileAvailability(filename string, response *Response) error {
	colorprint.Debug("INBOUND RPC REQUEST: Checking File Availability " + filename)
	fileSysLock.RLock()
	video, ok := localFileSys.Files[filename]
	if ok {
		colorprint.Info("File " + filename + " is available")
		colorprint.Debug("INBOUND RPC REQUEST COMPLETED")
		response.Avail = true
		response.SegNums = video.SegNums
		response.SegsAvail = video.SegsAvail
	} else {
		colorprint.Alert("File " + filename + " is unavailable")
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
func (service *Service) GetFileSegment(segReq *ReqStruct, segment *VidSegment) error {
	colorprint.Debug("INBOUND RPC REQUEST: Sending video segment for " + segReq.Filename)
	var seg VidSegment
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

		return errors.New("File unavailable.")
		fileSysLock.Unlock()
	}
	colorprint.Warning(outputstr)
	return nil
}
