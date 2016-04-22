package streamerServer

import (
	"net"
	//"os"
	"net/rpc"
	"log"
	"os/exec"
	"bytes"
	"fmt"
	//"strconv"
	"io/ioutil"
	"strings"
	//"../customChord"
)

type Reply struct {
  Val string
}
type Msg struct {
  Id int64
  Filename string
  Val string
  Address string
  Data []byte
}

var nodeAddr string
var nodeName string
var dest string
type NodeRPCService int

func (this *NodeRPCService) SaveToServer(msg *Msg, reply *Reply) error {
	fmt.Println("FILEFOLDERSHIT: ", msg.Filename)
	pathArr := strings.Split(msg.Filename, " ")

	path := "FFMPEG/NodesData/" + nodeName + "/" + pathArr[0] + "/" + pathArr[1]
	err := ioutil.WriteFile(path, msg.Data, 0644)
	checkError(err)
	reply.Val = "OKiE"

	// TODO: Need to update chord dataMap
	//customChord.SaveToStore(pathArr[0], pathArr[1], msg.Data)
	return nil
}


func (this *NodeRPCService) StartStreaming(msg *Msg, reply *Reply) error {
	// "udp://127.0.0.1:1234"
	//addr := "udp://127.0.0.1" + msg.Address

	var out bytes.Buffer
	var stderr bytes.Buffer
	

	// cmd0 := exec.Command("pwd")
	// cmd0.Stdout = &out
	// _ = cmd0.Start()
	// _ = cmd0.Wait()
	// fmt.Println("PWD: ", out.String())
	fnArr := strings.Split(msg.Filename, ".")
	path := dest + fnArr[0] + "/%05d.png"
	cmd := exec.Command("ffmpeg", "-start_number", msg.Val, "-re", "-i", path, "-r", "10", 
	"-vcodec", "mpeg4", "-f", "mpegts", msg.Address)

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()

	if err != nil {
    	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    	return err
	}
	//fmt.Println("Result: " + out.String())

	//checkError(err)
	log.Printf("Waiting to start streaming frames...")
	err = cmd.Wait()
	//fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	log.Printf("Frame streaming finished with error: %v", err)
	if err != nil {
    	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    	return err
	}
	reply.Val = "ok"
	return nil
}

/* 
* Set up the listener for RPC requests, serve the connections when required.
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
    go rpc.ServeConn(newRPCConnection) // Serve a request concurrently
  }
  rpcListener.Close()
}

func GetFrames(filename string) int64 {
	// ffmpeg -i sample.mp4 -r 100 -f image2 output/%05d.png
	fnArr := strings.Split(filename, ".")
	// destPath := "FFMPEG/NodesData/" + nodeName + "/output/%05d.png"
	destPath := "FFMPEG/NodesData/" + nodeName + "/" + fnArr[0] + "/%05d.png"
	sourcePath := "FFMPEG/NodesData/" + nodeName + "/source/" + filename

	// split video into png frames
	cmd := exec.Command("ffmpeg", "-i", sourcePath, "-r", "100", "-f",
		"image2", destPath)
	err := cmd.Start()
	checkError(err)
	log.Printf("Waiting for video to finish processing into individual frames...")
	err = cmd.Wait()
	log.Printf("Frame processing finished with error: %v", err)

	path := "FFMPEG/NodesData/" + nodeName + "/" + fnArr[0] + "/"
	files,_ := ioutil.ReadDir(path)
	numFrames := int64(len(files))

    return numFrames-1
}

func Start(rpcServerAddr string, name string) {
	nodeAddr = rpcServerAddr
	nodeName = name

	//createDirectories() TODO

	dest = "FFMPEG/NodesData/" + nodeName + "/"
	//getFrames(dest)
	log.Println("Launching rpc service to serve stream requests...")
	launchRPCService(nodeAddr)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}