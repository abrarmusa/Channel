package main

/*
  UBC CS416 Distributed Systems Project Source Code

  @author: Abrar, Ito, Mimi, Shariq
  @date: Mar. 1 2016 - Apr. 11 2016.

  Usage:
    go run node.go [node ip:port] [starter-node ip:port] [-r=replicationFactor] [-t]

    [node ip:port] : this node's ip/port combo
    [starter-node ip:port] : the entry point node's ip/port combo
    [-r=replicationFactor] : replication factor for keys, default r = 2
    [-t] : trace mode for debugging

  Copy/paste for quick testing:
    "go run node.go :0 :6666" <-- trace off, default replication factor (2)
    "go run node.go :0 :6666 -t -r=5" <-- trace on, r = 5
*/

import (
  "crypto/sha1"
  "encoding/hex"
  "encoding/json"
  "flag"
  "fmt"
  "net"
  "os"
  "time"
  "sync"
  "transfile"
)

// =======================================================================
// ====================== Global variables/types =========================
// =======================================================================
// Some types used for RPC communication.
type nodeRPCService int
type NodeMessage struct {
  Msg string
}
type ConnResponse struct {
  Conn *net.TCPConn
}

// Some static command line options.
var traceMode bool
var replicationFactor int

// The finger table.
// Reading? Wrap RLock/RUnlock: lock.RLock(), lock.fingerTable[..], lock.RUnlock()
// Writing? Wrap Lock/Unlock: lock.Lock(), lock.fingerTable[..] = .. ,lock.Unlock()
var lock = struct {
  sync.RWMutex
  fingerTable map[string]string
}{fingerTable: make(map[string]string)} // initialize it here (or in main)

// =======================================================================
// ======================= Function definitions ==========================
// =======================================================================
/* Updates or inserts a finger table entry at this node. Called by other nodes.
 */
func (this *nodeRPCService) UpdateFingerTableEntry(id string, addr string,
  reply *NodeMessage) error {
  // Lock 'em up.
  lock.Lock()
  defer lock.Unlock()

  // Update.
  lock.fingerTable[id] = addr

  // Let the other guy know it went well.
  reply.Msg = "Ok"
  return nil
}

/* Grabs a finger table entry at this node. Called by other nodes.
 */
func (this *nodeRPCService) GetFingerTableEntry(id string,
  reply *NodeMessage) error {
  // Lock 'em up (read mode).
  lock.RLock()
  defer lock.RUnlock()

  // Send the value of the entry that some guy requested.
  reply.Msg = lock.fingerTable[id]
  return nil
}

/* Removes a finger table entry at this node. Called by other nodes.
 */
func (this *nodeRPCService) DeleteFingerTableEntry(id string,
  reply *NodeMessage) error {
  // Lock 'em up.
  lock.Lock()
  defer lock.Unlock()

  // Update.
  delete(lock.fingerTable, id)

  // Let the other guy know it went well.
  reply.Msg = "Ok"
  return nil
}

/* Checks error value and prints/exits if non nil.
 */
func checkError(err error) {
  if err != nil {
    fmt.Println("Error string: ", err)
    os.Exit(-1)
  }
}

/* Returns the SHA1 hash value as a string, of a key k.
 */
func computeSHA1Hash(key string) string {
  buf := []byte(key)
  h := sha1.New()
  h.Write(buf)
  str := hex.EncodeToString(h.Sum(nil))
  return str
}

/* Computes the distance between two SHA1 hashes.
 */
func computeDistBetweenTwoHashes(key1 string, key2 string) int64 {
  // todo
  return 69
}

/* Send periodic heartbeats to let predecessor know this node is still alive.
 */
func sendAliveMessage(conn *net.TCPConn, id string) {
  for {
    // send this node's id as an alive message
    msg := NodeMessage{id}
    aliveMessage, err := json.Marshal(msg)
    checkError(err)
    b := []byte(aliveMessage)
    _, err = conn.Write(b)
    checkError(err)

    time.Sleep(5 * time.Second)
  }
}

/* Perform recursive search through finger tables to place me at the right spot.
 */
func locateSuccessor(conn *net.TCPConn, id string) *net.TCPConn {
  // recursive search through finger tables
  // use computeDistBetweenTwoHashes(key1 string, key2 string)

  // Send a special message of value "where", so that node knows it wants to find its place in the identifier circle.
  msg := NodeMessage{"where"}
  msgInJSON, err := json.Marshal(msg)
  checkError(err)
  buf := []byte(msgInJSON)
  // The listening function should have an interface{} to deal with a "where" message.
  // Handle this write in the listenForControlMessages function
  _, err = conn.Write(buf)
  checkError(err)

  // The response should hold the conn of the successor.
  n, err := conn.Read(buf)
  checkError(err)
  var successor ConnResponse
  err = json.Unmarshal(buf[:n], &successor)
  checkError(err)

  return successor.Conn
}

/* Handle messages from .
 */
func listenForControlMessages(nodeAddr string) {
  // Todo
}

/* Only called when there are no nodes yet that are started up. This node becomes the first node.
 */
func startUpSystem(nodeAddr string) {
  // Todo
}

/* Attempt to join the system given the ip:port of a running node.
 */
func connectToSystem(nodeAddr string, startAddr string) {
  // Get this node's IP hash, which will be used as its ID.
  id := computeSHA1Hash(nodeAddr)

  nodeTCPAddr, err := net.ResolveTCPAddr("tcp", nodeAddr)
  checkError(err)
  startTCPAddr, err := net.ResolveTCPAddr("tcp", startAddr)
  checkError(err)

  // Figure out where I am in the identifier circle.
  conn, err := net.DialTCP("tcp", nodeTCPAddr, startTCPAddr)
  successorConn := locateSuccessor(conn, id)

  // Don't need this conn object anymore.
  if successorConn != conn {
    defer conn.Close()
  }

  // Send a heartbeat every 5 secs. The successor will start tracking this
  // node if it hadn't sent an alive message before.
  go sendAliveMessage(successorConn, id)

  // Listen for TCP message requests to this node.
  listenForControlMessages(nodeAddr)
}

/* The main function.
 */
func main() {
  // Handle the command line.
  if len(os.Args) > 5 || len(os.Args) < 3 {
    fmt.Println("Usage: go run node.go [node ip:port] [starter-node ip:port] [-r=replicationFactor] [-t]")
    os.Exit(-1)
  } else {
    nodeAddr := os.Args[1] // ip:port of this node
    startAddr := os.Args[2] // ip:port of initial node
    flag.IntVar(&replicationFactor, "r", 2, "replication factor")
    flag.BoolVar(&traceMode, "t", false, "trace mode")
    flag.Parse()

    // Set both args as equal on the command line if no node is operational yet.
    // Need to start up a source node (this) first.
    if (nodeAddr == startAddr) {
      startUpSystem(nodeAddr)
    // Else join an existing identifier circle, given the address of a node that is running.
    } else {
      connectToSystem(nodeAddr, startAddr)
    }
  }

  // Setup UDP server
  u := transfile.UdpInfo{
    Conn: nil,
    Port: ":0",
  }
  u.SetupUdp()
  fmt.Println("This UDP:", u)

  /*
  (1) New u.Port is to be advertised as part of finger
  (2) Send data like this: u.SendUdpCall("198.162.33.54:40465")
  */

	go u.ReceiveUdpCall() // infinite waiting & receiving
	transfile.StayLive()
}
