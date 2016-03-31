package main
/*
  UBC CS416 Distributed Systems Project Source Code

  @author: Abrar, Ito, Mimi, Shariq
  @date: Mar. 1 2016 - Apr. 11 2016.

  Usage:
    go run node.go [node ip:port] [starter-node ip:port]"

    [node ip:port] : this node's ip/port combo
    [starter-node ip:port] : the entry point node's ip/port combo

  Copy/paste for quick testing:
    "go run node.go :6666 :6666" <- start up system at :6666
    "go run node.go :0 :6666" <- connect to node listening at :6666
*/

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net"
  "net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

// =======================================================================
// ====================== Global Variables/Types =========================
// =======================================================================
// Some types used for RPC communication.
type nodeRPCService int
type CommandMessage struct {
	Cmd        string
	SourceAddr string
	DestAddr   string
	Key        string
	Val        string
}

// Some static command line options.
var replicationFactor int
var m float64
var successor int64
var successorAddr string
var identifier int64

// The finger table.
var fingerLocker = struct {
	sync.RWMutex
	fingerTable map[int64]string
}{fingerTable: make(map[int64]string)}

// File storage table.
var fileLocker = struct {
	sync.RWMutex
	fileTable map[string]string
}{fileTable: make(map[string]string)}

// =======================================================================
// ======================= RPC Methods ===================================
// =======================================================================
/* Updates a finger table entry at this node. Called by other nodes.
 */
func (this *nodeRPCService) UpdateFingerTableEntry(id string, addr string,
  reply *CommandMessage) error {
  // Lock 'em up.
  fingerLocker.Lock()
  defer fingerLocker.Unlock()

  // Update.
  fingerLocker.fingerTable[id] = addr

  // Let the other guy know it went well.
  reply.Val = "Ok"
  return nil
}

// =======================================================================
// ======================= Helper Methods ================================
// =======================================================================
/* Checks error Value and prints/exits if non nil.
 */
func checkError(err error) {
	if err != nil {
		fmt.Println("Error string: ", err)
		os.Exit(-1)
	}
}

/* Returns the SHA1 hash Value as a string, of a Key k.
 */
func computeSHA1Hash(Key string) string {
	buf := []byte(Key)
	h := sha1.New()
	h.Write(buf)
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

// =======================================================================
// =================== Maintenance Methods ===============================
// =======================================================================
/* Send periodic heartbeats to let predecessor know this node is still alive.
 */
func sendAliveMessage(conn *net.UDPConn, src string, dest string) {
	for {
		// send this node's id as an alive message
		msg := CommandMessage{"alive", src, dest, "", ""}
		aliveMessage, err := json.Marshal(msg)
		checkError(err)
		b := []byte(aliveMessage)
		_, err = conn.Write(b)
		checkError(err)

		time.Sleep(5 * time.Second)
	}
}

// =======================================================================
// ======================= Logical Methods ===============================
// =======================================================================
/* Perform recursive search through finger tables to place me at the right spot.
 */
func locateSuccessor(conn *net.UDPConn, addr string) {
	fmt.Println("## Locating successor...")

	// Send a special discover message.
	msg := CommandMessage{"_discover", addr, "", "", ""}
	msgInJSON, err := json.Marshal(msg)
	checkError(err)

	buf := []byte(msgInJSON)
	_, err = conn.Write(buf)
	fmt.Println("## Sent command: ", string(buf[:]))
	checkError(err)
}

/* Computes an identifier based on 2^m value and SHA1 hash.
 */
func getIdentifier(key string) int64 {
	id := computeSHA1Hash(key)
	k := big.NewInt(0)
	if _, ok := k.SetString(id, 16); ok {
		fmt.Println("Number: ", k)
	} else {
		fmt.Println("Unable to parse into big int")
	}
	power := int64(math.Pow(2, m))
	ret := (k.Mod(k, big.NewInt(power))).Int64()
	fmt.Println("## Identifier is: ", ret)
	return ret
}

/* Retrieves a value from the finger table.
 */
func getVal(key string) (string, bool) {
	id := getIdentifier(key)
	fingerLocker.RLock()
	val := fingerLocker.fingerTable[id]
	fingerLocker.RUnlock()
	if val == "" {
		return val, false
	} else {
		return val, true
	}
}

/* Traverses finger table and finds the closest id to contact.
 */
func sendToNextBestNode(msg CommandMessage) {
	KeyIdentifier := getIdentifier(msg.Key)
	// find node in finger table which is closest to requested Key
	var closestNode string
	minDistanceSoFar := math.Pow(2.0, float64(m))

	fingerLocker.RLock()
	for nodeIden, nodeAddr := range fingerLocker.fingerTable {
		if nodeIden-KeyIdentifier < int64(minDistanceSoFar) {
			minDistanceSoFar = float64(nodeIden - KeyIdentifier)
			closestNode = nodeAddr
		}
	}
	fingerLocker.RUnlock()

	// send message to closestNode
	jsonMsg, err := json.Marshal(msg)
	checkError(err)
	buf := []byte(jsonMsg)
	sendMessage(closestNode, buf)
}

/* Sends a UDP message.
 */
func sendMessage(addr string, msg []byte) {
	fmt.Println("## Dialing to send message...")
	fmt.Println("## Address to dial: ", addr)
	fmt.Println("## Sending Message: ", string(msg))
	conn, err := net.Dial("udp", addr)
	checkError(err)
	_, err = conn.Write(msg)
	checkError(err)
}

/* Handle incoming udp writes.
 */
func listenForUDPMsgs(nodeAddr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", nodeAddr)
	checkError(err)
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	defer conn.Close()

	var msg CommandMessage
	buf := make([]byte, 2048)
	for {
		fmt.Println("## Waiting for packet to arrive on udp port...")
		n, _, err := conn.ReadFromUDP(buf)
		fmt.Println("## Received Command: ", string(buf[:n]))
		checkError(err)
		err = json.Unmarshal(buf[:n], &msg)
		checkError(err)
		fmt.Println("Cmd: ", msg.Cmd)

		switch msg.Cmd {
		case "_getVal":
			v, haveKey := getVal(msg.Key)
			if haveKey {
				// respond with Value
				responseMsg := CommandMessage{"_resVal", nodeAddr, msg.SourceAddr, msg.Key, v}
				resp, err := json.Marshal(responseMsg)
				checkError(err)
				buf = []byte(resp)
				// connect to source of request and send Value
				sendMessage(msg.SourceAddr, buf)
			} else {
				// send to next best node
				sendToNextBestNode(msg)
			}
		case "_setVal":
			_, haveKey := getVal(msg.Key)
			if haveKey {
				// change Value
				fileLocker.Lock()
				fileLocker.fileTable[msg.Key] = msg.Val
				fileLocker.Unlock()
				responseMsg := CommandMessage{"_resGen", nodeAddr, msg.SourceAddr, "", "Key Updated"}
				resp, err := json.Marshal(responseMsg)
				checkError(err)
				buf = []byte(resp)
				// connect to source of request and send Value
				sendMessage(msg.SourceAddr, buf)
			} else {
				// send to next best node
				sendToNextBestNode(msg)
			}
		case "_resDisc":
			successor = getIdentifier(msg.Val)
			successorAddr = msg.Val
			fmt.Println("## Successor updated to address: ", msg.Val)
			fmt.Println("## Successor Identifier is: ", successor)
		case "_discover":
			nodeIdentifier := getIdentifier(msg.SourceAddr)
			if successor == -1 {
				fmt.Println("## No successor in network. Setting now to new node...")
				fingerLocker.Lock()
				fingerLocker.fingerTable[nodeIdentifier] = msg.SourceAddr
				fingerLocker.Unlock()

				// notify new node of its successor (current successor)
				responseMsg := CommandMessage{"_resDisc", nodeAddr, msg.SourceAddr, "", nodeAddr}
				resMsg, err := json.Marshal(responseMsg)
				checkError(err)
				buf := []byte(resMsg)
				//fmt.Println("HERE 0")
				sendMessage(msg.SourceAddr, buf)
				//fmt.Println("HERE 1")
				// update successor to new node
				successor = nodeIdentifier
				successorAddr = msg.SourceAddr
				break
			}
			if nodeIdentifier >= identifier && nodeIdentifier <= successor {
				// incoming node belongs between this node and its current successor
				// update finger table
				fmt.Println("New node fits between me and my successor. Updating finger table...")
				fingerLocker.Lock()
				fingerLocker.fingerTable[nodeIdentifier] = msg.SourceAddr
				fingerLocker.Unlock()

				// notify new node of its successor (current successor)
				responseMsg := CommandMessage{"_resDisc", nodeAddr, msg.SourceAddr, "", successorAddr}
				resMsg, err := json.Marshal(responseMsg)
				checkError(err)
				buf := []byte(resMsg)
				sendMessage(msg.SourceAddr, buf)
				// update successor to new node
				successor = nodeIdentifier
				successorAddr = msg.SourceAddr
				break
			} else {
				// forward command to next best node
				sendToNextBestNode(msg)
				break
			}
		case "_keepalive":
			//sendAliveMessage(conn *net.UDPConn, addr string)
		}
	}
}

/* Attempt to join the system given the ip:port of a running node.
 */
func connectToSystem(nodeAddr string, startAddr string) {
	fmt.Println("## Connecting to system")

	nodeUDPAddr, err := net.ResolveUDPAddr("udp", nodeAddr)
	checkError(err)
	startUDPAddr, err := net.ResolveUDPAddr("udp", startAddr)
	checkError(err)
	conn, err := net.DialUDP("udp", nodeUDPAddr, startUDPAddr)
	checkError(err)
	defer conn.Close()

	// Figure out where I am in the identifier circle.
	locateSuccessor(conn, nodeAddr)
}

/* The main function.
 */
func main() {
	// Handle the command line.
	if len(os.Args) > 4 || len(os.Args) < 2 {
		fmt.Println("Usage: go run node.go [node ip:port] [starter-node ip:port] [m size] [-r=replicationFactor] [-t]")
		os.Exit(-1)
	} else {
		nodeAddr := os.Args[1]                    // ip:port of this node
		startAddr := os.Args[2]                   // ip:port of initial node
    replicationFactor = 2 // # replications
    m = 3 // size of identifier circle
		successor = -1

		if nodeAddr == startAddr {
			fmt.Println("## Booting up system at addr ", nodeAddr)
			listenForUDPMsgs(nodeAddr)
		} else {
			fmt.Println("## Attempting to connect to node at ", startAddr)
			connectToSystem(nodeAddr, startAddr)
			listenForUDPMsgs(nodeAddr)
		}
	}
}
