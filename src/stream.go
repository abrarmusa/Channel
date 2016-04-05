// REQUIRES A LOT OF REFACTORIZATION AND CLEANUP


package main

import (
	"net"
	"net/rpc"
	"log"
	"os/exec"
	"os"
	"fmt"
)
type Reply struct {
  Val string
}
type Msg struct {
  Id int64
  Val string
}
var nodeAddr string
var nodeName string
var dest string
type NodeRPCService int

func (this *NodeRPCService) StartStreaming(msg *Msg, reply *Reply) error {
	cmd := exec.Command("ffmpeg", "-start_number", msg.Val, "-re", "-i", dest, "-r", "10", 
	"-vcodec", "mpeg4", "-f", "mpegts",  "udp://127.0.0.1:1234")
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting to start streaming frames...")
	err = cmd.Wait()
	log.Printf("Frame streaming finished with error: %v", err)
	reply.Val = "ok"
	return nil
}

// func (this *NodeRPCService) StartStreaming1(msg *Msg, reply *Reply) error {
// 	cmd := exec.Command("ffmpeg", "-start_number", msg.Val, "-re", "-i", dest, "-r", "10", 
// 	"-vcodec", "mpeg4", "-f", "mpegts",  "udp://127.0.0.1:1234")
// 	err := cmd.Start()
// 	checkError(err)
// 	log.Printf("Waiting to start streaming frames...")
// 	err = cmd.Wait()
// 	log.Printf("Frame streaming finished with error: %v", err)
// 	reply.Val = "ok"
// 	return nil
// }


// func (this *NodeRPCService) UpdateFingerTableEntry(msg *Msg, reply *Reply) error {
//   reply.Val = ""
//   // Lock 'em up.
//   fingerLocker.Lock()
//   defer fingerLocker.Unlock()

//   // Update.
//   fingerLocker.fingerTable[msg.Id] = msg.Val

//   // Let the other guy know it went well.
//   reply.Val = "ok"
//   return nil
// }

/* Set up the listener for RPC requests, serve the connections when required.
 */
func launchRPCService(addr string) {
  // Set up RPC service
  server := new(NodeRPCService)
  rpc.Register(server)
  rpcAddr, err := net.ResolveTCPAddr("tcp", addr)
  checkError(err)
  rpcListener, err := net.ListenTCP("tcp", rpcAddr)
  checkError(err)

  // Listen for RPC requests and serve concurrently
  for {
    newRPCConnection, err := rpcListener.AcceptTCP()
    checkError(err)
    go rpc.ServeConn(newRPCConnection) // Serve a request in parallel
  }
  rpcListener.Close()
}

func getFrames(dest string) {
	// ffmpeg -i sample.mp4 -r 100 -f image2 output/%05d.png
	cmd := exec.Command("ffmpeg", "-i", "FFMPEG/sample.mp4", "-r", "100", "-f",
		"image2", dest)
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting for video to finish processing into individual frames...")
	err = cmd.Wait()
	log.Printf("Frame processing finished with error: %v", err)
}

func startStream(dest string) {
	// ffmpeg -re -i FFMPEG/output/output%05d.png -r 10 -vcodec mpeg4 -f mpegts udp://127.0.0.1:1234
	cmd := exec.Command("ffmpeg", "-re", "-i", dest, "-r", "10", 
		"-vcodec", "mpeg4", "-f", "mpegts",  "udp://127.0.0.1:1234")
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting to start streaming frames...")
	err = cmd.Wait()
	log.Printf("Frame streaming finished with error: %v", err)
	
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run stream.go [node ip:port] [node name]")
		os.Exit(-1)
	} else {
		nodeAddr = os.Args[1]
		nodeName = os.Args[2]
	}

	dest = "FFMPEG/NodesData/" + nodeName + "/output/%05d.png"
	//getFrames(dest)

	launchRPCService(nodeAddr)
	//startStream()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}