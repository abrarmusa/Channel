package main

/*
  UBC CS416 Distributed Systems Project Source Code

  @author: Abrar, Ito, Mimi, Shariq
  @date: Mar. 1 2016 - Apr. 11 2016.

  Usage:
    go run node.go [node ip:port] [starter-node ip:port] [-r=replicationFactor] [-t]

    [node ip:port] : this node's ip/port combo
    [starter-node ip:port] : the entry point node's ip/port combo
    [-r=replicationFactor] : replication factor for Keys, default r = 2
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
  "math"
  //"strconv"
  "math/big"
  //"transfile"
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
type CommandMessage struct {
  Cmd string
  SourceAddr string
  DestAddr string
  Key string
  Val string
}

// Some static command line options.
var traceMode bool
var replicationFactor int
var store map[string]string
var ftab map[int64]string
var successor int64
var successorAddr string
var identifier int64
var m float64

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

  // Send the Value of the entry that some guy requested.
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

/* Computes the distance between two SHA1 hashes.
 */
func computeDistBetweenTwoHashes(Key1 string, Key2 string) int64 {
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
func locateSuccessor(conn *net.UDPConn, id string) {
  // recursive search through finger tables
  // use computeDistBetweenTwoHashes(Key1 string, Key2 string)
  fmt.Println("Locating successor...")

  // Send a special message of Value "where", so that node knows it wants to find its place in the identifier circle.
  msg := CommandMessage{"_discover", id, "", "", ""}
  msgInJSON, err := json.Marshal(msg)
  checkError(err)
  fmt.Println("JSON: ", string(msgInJSON))
  buf := []byte(msgInJSON)
  // The listening function should have an interface{} to deal with a "where" message.
  // Handle this write in the listenForControlMessages function
  _, err = conn.Write(buf)
  fmt.Println("Sent command: ", string(buf[:]))
  checkError(err)

  // The response should hold the conn of the successor.
  // n, err := conn.Read(buf)
  // checkError(err)
  // var successor ConnResponse
  // err = json.Unmarshal(buf[:n], &successor)
  // checkError(err)

  // return successor.Conn
}

func getIdentifier(Key string) int64 {
  id := computeSHA1Hash(Key)
  k := big.NewInt(0)
  if _, ok := k.SetString(id, 16); ok {
    fmt.Println("Number: ", k)
  } else {
    fmt.Println("Unable to parse into big int")
  }
  //k, err := strconv.ParseInt(id, 16, 64)
  //checkError(err)
  power := int64(math.Pow(2, m))
  ret := (k.Mod(k, big.NewInt(power))).Int64()
  fmt.Println("Identifier is: ", ret)
  return ret
}

func getVal(Key string) (string, bool) {
  v := ftab[getIdentifier(Key)]
  if v == "" {
    return v, false
  } else {
    return v, true
  }
}

/* Handle messages from .
 */
func listenForControlMessages(nodeAddr string) {
  // Todo
}

func sendToNextBestNode(msg CommandMessage) {
  KeyIdentifier := getIdentifier(msg.SourceAddr)
  // find node in finger table which is closest to requested Key
  var closestNode string
  minDistanceSoFar := math.Pow(2.0, float64(m))
  for nodeIden, nodeAddr := range ftab {
    if nodeIden - KeyIdentifier < int64(minDistanceSoFar) {
      minDistanceSoFar = float64(nodeIden - KeyIdentifier)
      closestNode = nodeAddr
    }
  }
  // send message to closestNode
  jsonMsg, err := json.Marshal(msg)
  checkError(err)
  buf := []byte(jsonMsg)
  sendMessage(closestNode, buf)
  // fmt.Println("Dialing to next best node...")
  // conn, _ := net.Dial("udp", closestNode)
  // _, err = conn.Write(buf)
  // checkError(err)
  
}

func sendMessage(addr string, msg []byte) {
  //tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
  //checkError(err)
  fmt.Println("Dialing to send message...")
  fmt.Println("Address to dial: ", addr)
  fmt.Println("Sending Message: ", string(msg))
  conn, err := net.Dial("udp", addr)
  checkError(err)
  _, err = conn.Write(msg)
  checkError(err)
}

/* Only called when there are no nodes yet that are started up. This node becomes the first node.
 */
func startUpSystem(nodeAddr string) {
  // listen and conect to node

  serverAddr, err := net.ResolveUDPAddr("udp", nodeAddr)
  checkError(err)

  conn, err := net.ListenUDP("udp", serverAddr)
  checkError(err)
  defer conn.Close()

  var msg CommandMessage
  //jsonMsg, err := json.Marshal(msg)
  //checkError(err)
  buf := make([]byte, 2048)

  for {
    fmt.Println("Waiting for packet to arrive on udp port...")
    n, _, err := conn.ReadFromUDP(buf)
    fmt.Println("Received Command: ", string(buf[:n]))
    checkError(err)
    err = json.Unmarshal(buf[:n], &msg)
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
          store[msg.Key] = msg.Val
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
      case "_resDisc" :
        successor = getIdentifier(msg.Val)
        successorAddr = msg.Val
        fmt.Println("Successor updated to address: ", msg.Val)
        fmt.Println("Successor Identifier is: ", successor)
      case "_discover":
        nodeIdentifier := getIdentifier(msg.SourceAddr)
        if successor == -1 {
          fmt.Println("No successor in network. Setting now to new node...")
          ftab[nodeIdentifier] = msg.SourceAddr
          
          // notify new node of its successor (current successor)
          responseMsg := CommandMessage {"_resDisc", nodeAddr, msg.SourceAddr, "", nodeAddr}
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
          ftab[nodeIdentifier] = msg.SourceAddr
          // notify new node of its successor (current successor)
          responseMsg := CommandMessage {"_resDisc", nodeAddr, msg.SourceAddr, "", successorAddr}
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
    }
  }
}

/* Attempt to join the system given the ip:port of a running node.
 */
func connectToSystem(nodeAddr string, startAddr string) {
  // Get this node's IP hash, which will be used as its ID.
  //id := computeSHA1Hash(nodeAddr)
  fmt.Println("Connecting to peer system...")

  nodeUDPAddr, err := net.ResolveUDPAddr("udp", nodeAddr)
  checkError(err)
  startUDPAddr, err := net.ResolveUDPAddr("udp", startAddr)
  checkError(err)

  // Figure out where I am in the identifier circle.
  conn, err := net.DialUDP("udp", nodeUDPAddr, startUDPAddr)
  checkError(err)

  defer conn.Close()
  //successorConn := locateSuccessor(conn, nodeAddr)
  locateSuccessor(conn, nodeAddr)
  // Don't need this conn object anymore.
  // if successorConn != conn {
  //   defer conn.Close()
  // }

  // Send a heartbeat every 5 secs. The successor will start tracking this
  // node if it hadn't sent an alive message before.
  //go sendAliveMessage(successorConn, nodeAddr)

  // Listen for TCP message requests to this node.
  //listenForControlMessages(nodeAddr)
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

    store = make(map[string]string)
    ftab = make(map[int64]string)
    m = 7
    successor = -1
    successorAddr = ""

    fmt.Println("THIS NODE'S IDENTIFIER IS: ", getIdentifier(nodeAddr))

    // Set both args as equal on the command line if no node is operational yet.
    // Need to start up a source node (this) first.
    // if (nodeAddr == startAddr) {
    //   startUpSystem(nodeAddr)
    // // Else join an existing identifier circle, given the address of a node that is running.
    // } else {
    //   connectToSystem(nodeAddr, startAddr)
    // }
    if (nodeAddr == startAddr) {
      fmt.Println("First node in system. Listening for incoming connections...")
      startUpSystem(nodeAddr)
    } else {
      go func() {
        go startUpSystem(nodeAddr)
        connectToSystem(nodeAddr, startAddr)
      }()
    }
  }
  for {}

  // Setup UDP server
 //  u := transfile.UdpInfo{
 //    Conn: nil,
 //    Port: ":0",
 //  }
 //  u.SetupUdp()
 //  fmt.Println("This UDP:", u)

  
 //  (1) New u.Port is to be advertised as part of finger
 //  (2) Send data like this: u.SendUdpCall("198.162.33.54:40465")
  

	// go u.ReceiveUdpCall() // infinite waiting & receiving
	// transfile.StayLive()
}
