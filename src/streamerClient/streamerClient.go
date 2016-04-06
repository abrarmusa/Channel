package streamerClient

import (
	"log"
	"os/exec"
	"net/rpc"
	"fmt"
	"os"
)

type StreamNode struct {
	Name string
	Address string
	startSeqNum string
	numFrames int64
	Next *StreamNode
}
type Reply struct {
  Val string
}
type Msg struct {
  Id int64
  Val string
  Address string
}
var nodeAddr string
type NodeRPCService int

func ListenForStream(addr string) {
	// ffplay udp://127.0.0.1:1234
	go func() {
		cmd := exec.Command("ffplay", addr)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting to get stream...")
		err = cmd.Wait()
		log.Printf("Getting stream finished with error: %v", err)
	}()
}

func GetRpcHandler(rpcAddr string) (*rpc.Client) {
	//var err error
	var nodeRPCHandler *rpc.Client
	nodeRPCHandler, err := rpc.Dial("tcp", rpcAddr)
	
	checkError(err)
	//defer nodeRPCHandler.Close()
	return nodeRPCHandler
}

func StartStreaming(handler *rpc.Client, iden int64, startFrame string, addr string) {
	msg := Msg {iden, startFrame, addr}
	var reply Reply
	err := handler.Call("NodeRPCService.StartStreaming", &msg, &reply) // returns id in msg.Id, and ip:port in msg.Val
	checkError(err)
	fmt.Println("Reply received: ", reply.Val)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(-1)
	}
}



























// =============================== PREVIOUS UNMERGED CODE =====================================//

// package main
// /*
//   UBC CS416 Distributed Systems Project Source Code

//   @author: Abrar, Ito, Mimi, Shariq
//   @date: Mar. 1 2016 - Apr. 11 2016.

//   Usage:
//     go run stream_client.go [starter-node ip:port]"

//     [client ip:port] : this client's ip:port combo
//     [starter-node ip:port] : the entry point node's ip/port combo

//   Copy/paste for quick testing:
//     "go run stream_client.go :1234 :6666" <- connect to node listening at :6666
// */

// import (
// <<<<<<< HEAD
// 	//"bytes"
// 	//"encoding/json"
// 	//"fmt"
// 	"log"
// 	"os/exec"
// )

// func getStream() {
// 	// ffplay udp://127.0.0.1:1234
// 	cmd := exec.Command("ffplay", "udp://127.0.0.1:1234")
// 	err := cmd.Start()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("Waiting to get stream...")
// 	err = cmd.Wait()
// 	log.Printf("Getting stream finished with error: %v", err)
// }
// =======
//   "encoding/json"
// 	"fmt"
//   "net"
//   "os"
//   vlc "github.com/adrg/libvlc-go"
// )

// // =======================================================================
// // ====================== Global Variables/Types =========================
// // =======================================================================
// type CommandMessage struct {
//   Cmd string
//   SourceAddr string
//   DestAddr string
//   Key string
//   Val string
//   Store map[string]string
// }

// var myAddr string
// var startAddr string

// // =======================================================================
// // ============================ Methods ==================================
// // =======================================================================
// /* Checks error Value and prints/exits if non nil.
//  */
// func checkError(err error) {
//   if err != nil {
//     fmt.Println("Error: ", err)
//     os.Exit(-1)
//   }
// }

// // Basic VLC setup
// // https://github.com/adrg/libvlc-go
// func initializeVLC() {
//   if err := vlc.Init(); err != nil {
//       fmt.Println(err)
//       return
//   }
//   defer vlc.Release()

//   // Create a new player
//   player, err := vlc.NewPlayer()
//   if err != nil {
//     fmt.Println(err)
//     return
//   }
//   defer func() {
//     player.Stop()
//     player.Release()
//   }()

//   // Set player media. The second parameter of the method specifies if
//   // the media resource is local or remote.
//   // err = player.SetMedia("localPath/test.mp4", true)
//   err = player.SetMedia("http://stream-uk1.radioparadise.com/mp3-32", false)
//   if err != nil {
//     fmt.Println(err)
//     return
//   }

//   // Play
//   err = player.Play()
//   if err != nil {
//     fmt.Println(err)
//     return
//   }
// }

// // Feeds a VLC player a stream of bytes
// func concatenateStream(msg *CommandMessage) {
// >>>>>>> 9725218a4dfaceb871b0715ee270dfb876ec472b

// func main() {
// 	getStream()
// }

// // Sends a file over to a running node.
// // Assumption: same computer.
// func uploadFile(filename string, filepath string) {
//   conn, err := net.Dial("udp", startAddr)
//   checkError(err)
//   defer conn.Close()

//   // Sends a special message to be handled by the node
//   msg := CommandMessage{"_upload", myAddr, startAddr, filename, filepath, nil}
//   msgInJSON, err := json.Marshal(msg)
//   checkError(err)
//   buf := []byte(msgInJSON)
//   _, err = conn.Write(buf)
//   checkError(err)

//   fmt.Println("Uploaded: ", filename)
// }

// // Opens vlc as bytes become available through a stream to connected node.
// func playFile(filename string) {
//   myUDPAddr, err := net.ResolveUDPAddr("udp", myAddr)

//   // Notify the node what file I want to stream.
//   conn, err := net.Dial("udp", startAddr)
//   checkError(err)
//   defer conn.Close()
//   sendMsg := CommandMessage{"_stream", myAddr, startAddr, filename, "", nil}
//   msgInJSON, err := json.Marshal(sendMsg)
//   checkError(err)
//   buf := []byte(msgInJSON)
//   _, err = conn.Write(buf)
//   checkError(err)

//   // Listen for a response and play accordingly.
//   conn2, err := net.ListenUDP("udp", myUDPAddr)
//   checkError(err)
//   defer conn2.Close()

//   initializeVLC()

//   // Feed to VLC as bytes are received.
//   var msg CommandMessage
//   buf = make([]byte, 1028)
//   for {
//     n, _, err := conn2.ReadFromUDP(buf)
//     fmt.Println("Received Response: ", string(buf[:n]))
//     checkError(err)
//     err = json.Unmarshal(buf[:n], &msg)

//     // Note done yet: add bytes to VLC player
//     if msg.Cmd != "EOF" {
//       concatenateStream(&msg)

//     // No bytes left to stream if Cmd == EOF
//     } else {
//       break
//     }
//   }
// }

// // Main f'n
// func main() {
//   if len(os.Args) != 3 {
//     fmt.Println("Usage: go run stream_client.go [client ip:port] [starter-node ip:port]")
//     os.Exit(-1)
//   } else {
//     myAddr = os.Args[1]
//     startAddr = os.Args[2]                   // ip:port of initial node
//     uploadFile("sample.mp4", "./")
//     playFile("sample.mp4")
//   }
// }
