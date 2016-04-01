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
	"fmt"
	"math"
	"math/big"
	"net"
  "net/rpc"
	"os"
  //"strconv"
	"sync"
	"time"
)

// =======================================================================
// ====================== Global Variables/Types =========================
// =======================================================================
// Some types used for RPC communication.
type Reply struct {
  Val string
}
type Msg struct {
  Id int64
  Val string
}
type NodeRPCService int

// Some static variables.
var myIdentifier int64
var replicationFactor int
var m float64

// The finger table.
var fingerLocker = struct {
	sync.RWMutex
	fingerTable map[int64]string
}{fingerTable: make(map[int64]string)}

// File storage table.
var fileLocker = struct {
	sync.RWMutex
	fileTable map[int64]string
}{fileTable: make(map[int64]string)}

// =======================================================================
// ======================= Helper Methods ================================
// =======================================================================
/* Checks error Value and prints/exits if non nil.
 */
func checkError(err error) {
  if err != nil {
    fmt.Println("Error: ", err)
    os.Exit(-1)
  }
}

/* Returns the SHA1 hash Value as a string, of a Key k.
 */
func computeSHA1Hash(id string) string {
  buf := []byte(id)
  h := sha1.New()
  h.Write(buf)
  str := hex.EncodeToString(h.Sum(nil))
  return str
}

/* Computes an identifier based on 2^m value and SHA1 hash.
 */
func computeIdentifier(id string) int64 {
  hash := computeSHA1Hash(id)
  k := big.NewInt(0)
  if _, ok := k.SetString(hash, 16); ok {
  } else {
    fmt.Println("--> Unable to parse into big int")
  }
  power := int64(math.Pow(2, m))
  ret := (k.Mod(k, big.NewInt(power))).Int64()
  fmt.Println("--> Computed identifier: ", ret)
  return ret
}

/* Prints the finger table entries to standard output.
 */
func printFingerTable() {
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf(" Finger table (unordered) for this node: %d\n", myIdentifier)
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf("| ID   |    VAL    |\n")
  fingerLocker.RLock()
  defer fingerLocker.RUnlock()

  // Runs up to size m.
  for id := range fingerLocker.fingerTable {
    fmt.Printf("| %3d  | %9s |\n", id, fingerLocker.fingerTable[id])
  }
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
}

/* Prints the finger table entries to standard output.
 */
func printFileTable() {
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf(" File table (unordered) for this node: %d\n", myIdentifier)
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf("|  ID  |     VAL    |\n")
  fileLocker.RLock()
  defer fileLocker.RUnlock()

  // Runs up to size m.
  for id := range fileLocker.fileTable {
    fmt.Printf("| %3d  | %10s |\n", id, fileLocker.fileTable[id])
  }
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
}

/* Returns true if an identifier is between this node and successor on the id circle.
 */
func isBetweenMeAndSuccessor(iden int64, me int64, suc int64) bool {
  if suc < me {
    if iden >= me && iden >= suc {
      return true
    } else if iden < me && iden <= suc {
      return true
    } else {
      return false
    }
  } else if suc >= me {
    if iden >= me && iden <= suc {
      return true
    } else {
      return false
    }
  }
}

// =======================================================================
// ======================= RPC Methods ===================================
// =======================================================================
/* Updates a finger table entry at this node. Called during entry/exit/failures.
 */
func (this *NodeRPCService) UpdateFingerTableEntry(msg *Msg, reply *Reply) error {
  // Lock 'em up.
  fingerLocker.Lock()
  defer fingerLocker.Unlock()

  // Update.
  fingerLocker.fingerTable[msg.Id] = msg.Val

  // Let the other guy know it went well.
  reply.Val = "Ok"
  return nil
}

/* Updates a file table entry at this node.
 */
func (this *NodeRPCService) UpdateFileTableEntry(msg *Msg, reply *Reply) error {
  // Lock 'em up.
  fileLocker.Lock()
  defer fileLocker.Unlock()

  // Update.
  fileLocker.fileTable[msg.Id] = msg.Val

  // Let the other guy know it went well.
  reply.Val = "Ok"
  return nil
}

/* Attempt to locate appropriate successor for caller node.
 */
func (this *NodeRPCService) Discover(msg *Msg, reply *Reply) error {
  // Lock for reading.
  fingerLocker.RLock()
  defer fingerLocker.RUnlock()

  // Find the best fit given an id.
  for id := range fingerLocker.fingerTable {
    if (id >= myIdentifier) && (id <= startIdentifier) {
      fingerLocker.fingerTable[id] = startAddr
    } else {
      fingerLocker.fingerTable[id] = nodeAddr
    }
  }

  // Let the other guy know it went well.
  reply.Val = "Ok"
  return nil
}

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

// =======================================================================
// ======================= Logical Methods ===============================
// =======================================================================
/* Sets up the finger table entries from 0 to 2^m.
 */
func initializeFingerTable() {
  fingerLocker.Lock()
  defer fingerLocker.Unlock()

  // (id + 2^i) mod m
  for i := 0; i < int(m); i++ {
    id := int64(math.Mod(float64(myIdentifier) + math.Pow(2, float64(i)),
                 math.Pow(2, float64(m))))
    fingerLocker.fingerTable[id] = "" // to be populated later
  }
}

/* Send periodic heartbeats to let predecessor know this node is still alive.
 */
func sendAliveMessage(conn *net.UDPConn, src string, dest string) {
  var reply Reply
  for {
    // send this node's id as an alive message
    reply.Val = ""
    aliveMessage, err := json.Marshal(reply)
    checkError(err)
    b := []byte(aliveMessage)
    _, err = conn.Write(b)
    checkError(err)

    time.Sleep(5 * time.Second)
  }
}

/* Attempt to join a system given the ip:port of a running node.
 */
func connectToSystem(nodeAddr string, startAddr string) {
  nodeRPCHandler, err := rpc.Dial("tcp", startAddr)
  checkError(err)
  defer nodeRPCHandler.Close()

  // Initialize finger table entries according to start node.
  startIdentifier := computeIdentifier(startAddr)
  fingerLocker.Lock()
  for id := range fingerLocker.fingerTable {
    if (id >= myIdentifier) && (id <= startIdentifier) {
      fingerLocker.fingerTable[id] = startAddr
    } else {
      fingerLocker.fingerTable[id] = nodeAddr
    }
  }
  fingerLocker.Unlock()
  printFingerTable()

  // Initialize values (base values of RPC chain)
  msg := Msg {myIdentifier, ""}
  currentClosestSuccessor := startIdentifier // recursively gets smaller
  currentClosestValue := startAddr
  var reply Reply

  // Start an RPC chain to place to populate finger table entries.
  for {
    err = nodeRPCHandler.Call("NodeRPCService.Discover", &msg, &reply) // returns id in msg.Id, and ip:port in msg.Val
    checkError(err)
    if isBetweenMeAndSuccessor(msg.Id, myIdentifier, currentClosestSuccessor)
        && (reply.Val == "ok") {
      currentClosestSuccessor = msg.Id
      currentClosestValue = msg.Val
    } else if (reply.Val != "ok")  {
      fmt.Println("--> Something went wrong with RPC.Discover: ", reply.Val)
      break
    } else {
      break
    }
  }

  // Re-examine the table.
  printFingerTable()
}


/* The main function.
 */
func main() {
	// Handle the command line.
	if len(os.Args) > 4 || len(os.Args) < 2 {
		fmt.Println("Usage: go run node.go [node ip:port] [starter-node ip:port]")
		os.Exit(-1)
	} else {
		nodeAddr := os.Args[1]                    // ip:port of this node
		startAddr := os.Args[2]                   // ip:port of initial node
    replicationFactor = 2 // size of replication window
    m = 7 // size of identifier circle
    myIdentifier = computeIdentifier(nodeAddr) // node's identifier based on ip:port and m

    // Initialize finger table entries from 0 to m for this node.
    initializeFingerTable()

    // Give the same nodeAddr and startAddr if there are no running nodes.
    if nodeAddr == startAddr {
			fmt.Println("--> Booting up system at addr ", nodeAddr)
		} else {
			fmt.Println("--> Attempting to connect to node at ", startAddr)
			connectToSystem(nodeAddr, startAddr)
		}

    // Called by other nodes for rpc methods on this node's finger table.
    launchRPCService(nodeAddr)
	}
}
