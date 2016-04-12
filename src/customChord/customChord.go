package customChord

/*
  UBC CS416 Distributed Systems Project Source Code

  @author: Abrar, Ito, Mimi, Shariq
  @date: Mar. 1 2016 - Apr. 11 2016.

  Usage:
    go run node.go [node ip:port] [starter-node ip:port]

    [node ip:port] : this node's ip/port combo
    [starter-node ip:port] : the entry point node's ip/port combo

  Copy/paste for quick testing:
    "go run node.go :0 :6666"
*/

import (
  "crypto/sha1"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "net"
  "io"
  "io/ioutil"
  "os"
  "time"
  "math"
  "strconv"
  "math/big"
  "runtime"
  //"strings"
  //"./lib/fileshare"
  "reflect"
)

// =======================================================================
// ====================== Global variables/types =========================
// =======================================================================

type CommandMessage struct {
  Cmd string
  SourceAddr string
  DestAddr string
  Key string
  Val string
  Store map[string]VidFrames
  Type string
}

type VidFrames struct {
  Name string
  FrameStart string
  TotalFrames int64
  Data map[string][]byte
}

var dataMap map[string]VidFrames
var dataMapSuccessor map[string]VidFrames
var dataMapPredecessor map[string]VidFrames

var traceMode bool
var replicationFactor int
var store map[string]string
var backupStorePred map[string]string
var backupStoreSuc map[string]string
var ftab map[int64]string
var successor int64
var successorAddr string
var predecessor int64
var predecessorAddr string
var identifier int64
var m float64
var c chan string
var myAddr string
var fileTransferAddr string

var streamServerAddress string
var streamClientAddress string
var nodename string

var successorAliveChannel chan bool
var predecessorAliveChannel chan bool
var streamServerChannel chan string
var fileTransferChannel chan string

// =======================================================================
// ======================= Function definitions ==========================
// =======================================================================

/*
* Checks error value and prints/exits if non nil.
*/
func checkError(err error) {
  if err != nil {
    // fmt.Println("Error string: ", err)
    os.Exit(-1)
  }
}

/*
* Returns the SHA1 hash Value of input key as a string
*/
func computeSHA1Hash(Key string) string {
  buf := []byte(Key)
  h := sha1.New()
  h.Write(buf)
  str := hex.EncodeToString(h.Sum(nil))
  return str
}

/* Prints the finger table entries to standard output.
 */
func printFingerTable() {
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf(" Finger table (unordered) for this node: %d\n", identifier)
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf("| ID   |    VAL    |\n")

  // Runs up to size m.
  for id := range ftab {
    fmt.Printf("| %3d  | %9s |\n", id, ftab[id])
  }
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
}

/*
* Send a heartbeat message to let inquiring node know that we're still alive
*/
func sendAliveMessage(addr string) {
      msg := CommandMessage{"_heartbeat", myAddr, addr, strconv.FormatInt(identifier, 10), myAddr, dataMap, ""}
      aliveMessage, err := json.Marshal(msg)
      checkError(err)
      b := []byte(aliveMessage)
      sendMessage(addr, b)
}

/*
* Ask a node if it is alive
*/
func askIfAlive(timeout chan bool, addr string) {
    msg := CommandMessage{"_alive?", myAddr, addr, strconv.FormatInt(identifier, 10), myAddr, nil, ""}
    aliveMessage, err := json.Marshal(msg)
    checkError(err)
    b := []byte(aliveMessage)
    sendMessage(addr, b)
    time.Sleep(5 * time.Second)
    timeout <- true
}

// TODO: USE time.After instead of time.Sleep

/*
* Heartbeat handlers which periodically send heartbeat messages to both successor and predecessor
*/
func handlePredecessorHeartbeats() {

  for {
    if predecessorAddr != "" {
      timeout_p := make(chan bool, 1)
      go askIfAlive(timeout_p, predecessorAddr)

      select {
      case <-predecessorAliveChannel:
        fmt.Println("Heartbeat from predecessor", predecessorAddr)
      case <-timeout_p:
        fmt.Println("Timed out on predecessor heartbeat")
        predecessor = -1
        predecessorAddr = ""
        //stabilizeNode()
      }
      time.Sleep(2 * time.Second)
    } else {
      //fmt.Println("Looping until we have both successor and predecessor")
      runtime.Gosched()
    }
  }
}

func handleSuccessorHeartbeats() {

  for {
    if successorAddr != "" {
      timeout_s := make(chan bool, 1)
      go askIfAlive(timeout_s, successorAddr)

      select {
      case <-successorAliveChannel:
        fmt.Println("Heartbeat from successor: ", successorAddr)
      case <-timeout_s:
        // timed out, my successor might be dead. time to make some changes in our secret circle
        fmt.Println("Timed out on successor heartbeat")
        // locate new successor if any
        // update predecessor of new
        successor = -1
        successorAddr = ""
        stabilizeNode("successor")
      }
      time.Sleep(2 * time.Second)
    } else {
      //fmt.Println("Looping until we have both successor and predecessor")
      runtime.Gosched()
    }
  }
}


// /*
// * Ask a node if it is alive
// */
// func askIfAlive(timeout chan bool, addr string, reltype string) {
//     msg := CommandMessage{"_alive?", myAddr, addr, strconv.FormatInt(identifier, 10), myAddr, nil, reltype}
//     aliveMessage, err := json.Marshal(msg)
//     checkError(err)
//     b := []byte(aliveMessage)
//     sendMessage(addr, b)
//     time.Sleep(10 * time.Second)
//     timeout <- true
// }

// // TODO: USE time.After instead of time.Sleep

// /*
// * Heartbeat handlers which periodically send heartbeat messages to both successor and predecessor
// */
// func handlePredecessorHeartbeats() {

//   for {
//     if predecessorAddr != "" {
//       timeout_p := make(chan bool, 1)
//       go askIfAlive(timeout_p, predecessorAddr, "predecessor")

//       select {
//       case <-predecessorAliveChannel:
//         fmt.Println("Heartbeat from predecessor", predecessorAddr)
//       case <-timeout_p:
//         fmt.Println("Timed out on predecessor heartbeat")
//         predecessor = -1
//         predecessorAddr = ""
//         //stabilizeNode()
//         // TODO: TRANSFER ALL PREDECESSOR FILES TO NEW PREDECESSOR
//       }
//       runtime.Gosched()
//       time.Sleep(10 * time.Second)
//     } else {
//       //fmt.Println("Looping until we have both successor and predecessor")
//       runtime.Gosched()
//     }
//   }
// }

// func handleSuccessorHeartbeats() {

//   for {
//     if successorAddr != "" {
//       timeout_s := make(chan bool, 1)
//       go askIfAlive(timeout_s, successorAddr, "successor")

//       select {
//       case <-successorAliveChannel:
//         fmt.Println("Heartbeat from successor: ", successorAddr)
//       case <-timeout_s:
//         // timed out, my successor might be dead. time to make some changes in our secret circle
//         fmt.Println("Timed out on successor heartbeat")
//         // locate new successor if any
//         // update predecessor of new
//         successor = -1
//         successorAddr = ""
//         stabilizeNode("successor")
//         // TODO: TRANSFER ALL SUCCESSOR FILES TO NEW SUCCESSOR
//       }
//       runtime.Gosched()
//       time.Sleep(5 * time.Second)
//     } else {
//       //fmt.Println("Looping until we have both successor and predecessor")
//       runtime.Gosched()
//     }
//   }
// }

/*
* Stabilizes a node by finding a new successor. Ran after successor node dies.
*/
func stabilizeNode(position string) {
  
  for _, addr := range ftab {
    msg := CommandMessage{"_proposal", myAddr, addr, position, strconv.FormatInt(identifier, 10), dataMapSuccessor, ""}
    buf := getJSONBytes(msg)
    sendMessage(addr, buf)
  }

  // also send to predecessor
  msg := CommandMessage{"_proposal", myAddr, predecessorAddr, position, strconv.FormatInt(identifier, 10), dataMapSuccessor, ""}
  buf := getJSONBytes(msg)
  sendMessage(predecessorAddr, buf)

}

/*
* Find this node's predecessor
*/
func locatePredecessor(conn net.Conn) {
  msg := CommandMessage{"_locPred", myAddr, "", strconv.FormatInt(identifier, 10), "", nil, ""}
  msgInJSON, err := json.Marshal(msg)
  checkError(err)
  buf := []byte(msgInJSON)
  _, err = conn.Write(buf)
  checkError(err)
}

/*
* Find this node's successor
*/
func locateSuccessor(conn net.Conn, id string) {
  msg := CommandMessage{"_discover", id, "", "", "", nil, ""}
  msgInJSON, err := json.Marshal(msg)
  checkError(err)
  buf := []byte(msgInJSON)
  _, err = conn.Write(buf)
  checkError(err)
}

/*
* Inquire a node about where the identifier iden should lie on the Identifier Circle
*/
func getNodeInfo(nodeAddr string, iden int64, forType string) {
  msg := CommandMessage{"_getInfo", nodeAddr, successorAddr, "", strconv.FormatInt(iden, 10), nil, forType}
  jsonMsg, err := json.Marshal(msg)
  checkError(err)
  b := []byte(jsonMsg)
  sendMessage(successorAddr, b)
}

/*
* Initializes finger table populating entries from iden+2^0 to iden+2^m
*/
func initFingerTable(conn net.Conn, nodeAddr string) {
  thisIden := GetIdentifier(nodeAddr)
  for i := 0; i < int(m); i++ {
    key := int64( math.Mod( float64(thisIden) + math.Pow(2, float64(i)), math.Pow(2, float64(m)) ) )
    getNodeInfo(nodeAddr, key, "ftab")
  }
}

/*
* Returns an int64 identifier for an input key
*/
func GetIdentifier(Key string) int64 {
  id := computeSHA1Hash(Key)
  k := big.NewInt(0)
  if _, ok := k.SetString(id, 16); ok {
    //fmt.Println("Number: ", k)
  } else {
    // fmt.Println("Unable to parse into big int")
  }
  power := int64(math.Pow(2, m))
  ret := (k.Mod(k, big.NewInt(power))).Int64()

  fmt.Println("Identifier for ", Key, " : ", ret)

  return ret
}

/*
* Returns address of node with hash Key by doing a lookup in the finger table
* Second return value is true if lookup is successful
* Else returns false as the second return value
*/
func getVal(Key string) (string, bool) {
  v := ftab[GetIdentifier(Key)]
  if v == "" {
    return v, false
  } else {
    return v, true
  }
}

/*
* Sends to next best candidate for finding KeyIdentifier by searching through finger table
*/
func sendToNextBestNode(KeyIdentifier int64, msg CommandMessage) {
  var closestNode string
  minDistanceSoFar := int64(math.MaxInt64)
  for nodeIden, nodeAddr := range ftab {
    diff := nodeIden - KeyIdentifier
    if diff < minDistanceSoFar {
      minDistanceSoFar = diff
      closestNode = nodeAddr
    }
  }
  // send message to closestNode
  jsonMsg, err := json.Marshal(msg)
  checkError(err)
  buf := []byte(jsonMsg)
  sendMessage(closestNode, buf)
}

/*
* Sends a message msg to node with address addr
*/
func sendMessage(addr string, msg []byte) {
  //fmt.Println("Dialing to send message...")
  //fmt.Println("Address to dial: ", addr)
  if addr == "" {
    // send to self(?) for now - testing streaming
    fmt.Println("Sending message to self...")
    sendMessage(myAddr, msg)
    return
  }
  //fmt.Println("Sending Message: ", string(msg))
  conn, err := net.Dial("udp", addr)
  checkError(err)

  defer conn.Close()
  _, err = conn.Write(msg)
  checkError(err)
}

/*
* Checks if an identifier iden lies between this node and its successor
*/
func betweenIdens(suc int64, me int64, iden int64) bool {
  if suc < me {
    if iden > me && iden > suc {
      return true
    } else if iden < me && iden < suc {
      return true
    }
  } else if suc > me {
    if iden > me && iden < suc {
      return true
    }
  }
  return false
}

/*
* Replies with information about node where the inquired identifier should belong
* If it can't, sends the message to next best node in finger table
*/
func provideInfo(msg CommandMessage, nodeAddr string) {
  iden, err := strconv.ParseInt(msg.Val, 10, 64)
  checkError(err)
  if betweenIdens(successor, identifier, iden) {
    reply := CommandMessage{"_resInfo", nodeAddr, msg.SourceAddr, msg.Val, successorAddr, nil, msg.Type}
    jsonReply, err := json.Marshal(reply)
    checkError(err)
    b := []byte(jsonReply)
    sendMessage(msg.SourceAddr, b)
  } else if identifier == iden {
    // heloo.. is it me you're looking for
    reply := CommandMessage{"_resInfo", nodeAddr, msg.SourceAddr, msg.Val, nodeAddr, nil, msg.Type}
    jsonReply, err := json.Marshal(reply)
    checkError(err)
    b := []byte(jsonReply)
    sendMessage(msg.SourceAddr, b)
  } else if val, ok := ftab[iden]; ok {
    reply := CommandMessage{"_resInfo", nodeAddr, msg.SourceAddr, msg.Val, val, nil, msg.Type}
    jsonReply, err := json.Marshal(reply)
    checkError(err)
    b := []byte(jsonReply)
    sendMessage(msg.SourceAddr, b)
  } else {
    // fmt.Println("Can't provide info, forwarding message to next best node")
    sendToNextBestNode(iden, msg)
  }
}

/*
* Sends a message with predecessor info
*/
func sendPredInfo(src string, succ string) {
  responseMsg := CommandMessage{"_resLocPred", myAddr, src, "predecessor", succ, nil, ""}
  resp, err := json.Marshal(responseMsg)
  checkError(err)
  buf := []byte(resp)
  sendMessage(src, buf)
}

/*
* Initializes the P2P system
* Responsible for triggering heartbeat goroutines, backup goroutine and command loop
*/
func startUpSystem(nodeAddr string) {

  serverAddr, err := net.ResolveUDPAddr("udp", nodeAddr)
  checkError(err)

  //go handlePredecessorHeartbeats()
  //go handleSuccessorHeartbeats()

  //go maintainBackup()


  // fmt.Println("Trying to listen on: ", serverAddr)
  conn, err := net.ListenUDP("udp", serverAddr)
  checkError(err)
  defer conn.Close()

  var msg CommandMessage
  buf := make([]byte, 2048)

  for {
    // fmt.Println("Waiting for packet to arrive on udp port...")
    n, _, err := conn.ReadFromUDP(buf)
    // fmt.Println("Received Command: ", string(buf[:n]))
    checkError(err)
    err = json.Unmarshal(buf[:n], &msg)
    k, err := strconv.ParseInt(msg.Key, 10, 64)
    //checkError(err)

    switch msg.Cmd {
      case "_copyFiles":
        fmt.Println("Going to copy all received files from: ", msg.SourceAddr)
        // copy files and adjust store (?)
        copyFiles(msg.Store)

      case "_fileProposal":
        fmt.Printf("Received proposal for file %s with identifier %s\n", msg.Key, msg.Val)
        cmdMsg := CommandMessage{"_resFileProposal", myAddr, msg.SourceAddr, msg.Key, fileTransferAddr, nil, msg.Type}
        b := getJSONBytes(cmdMsg)
        sendMessage(msg.SourceAddr, b)
      case "_resFileProposal":
        fmt.Println("Received response for file proposal: ", msg.Key)
        store[msg.Key] = "available"
        fileTransferChannel <- msg.Val
      case "_stream":
        fmt.Println("Received _stream command from: ", msg.SourceAddr)
        cmdMsg := CommandMessage{"_resStream", myAddr, msg.SourceAddr, "Stream Server Address", streamServerAddress, nil, msg.Type}
        b := getJSONBytes(cmdMsg)
        sendMessage(msg.SourceAddr, b)
      case "_resStream":
        fmt.Println("Received _resStream from: ", msg.SourceAddr)
        streamServerChannel <- msg.Val
      case "_storeBackup":
        sendKeyMap(msg.SourceAddr)
      case "_resStoreBackup":
        // REDUNDANT
        // fmt.Println("Setting backup store to: ", msg.Store)
        // backupStore = msg.Store
      case "_copyStore":
        // DONT KNOW IF I NEED THIS
        // fmt.Println("Received command to set our store to provided store: ", msg.Store)
        //store = msg.Store
      case "_proposal":
        if msg.Key == "successor" && predecessor == -1 {
          // accept proposal

          // send a message
          responseMsg := CommandMessage{"_resProposal", myAddr, msg.SourceAddr, "successor", strconv.FormatInt(identifier, 10), nil, ""}
          b := getJSONBytes(responseMsg)
          sendMessage(msg.SourceAddr, b)
        } else if predecessor != -1 && predecessorAddr != "" {
          // i have a predecessor, send message to my predecessor, passing along the chain till a node with no predecessor
          b := getJSONBytes(msg)
          sendMessage(predecessorAddr, b)
        }
      case "_resProposal":
        if msg.Key == "successor" && successor == -1 {
          successor, _ = strconv.ParseInt(msg.Val, 10, 64)
          successorAddr = msg.SourceAddr
          initFingerTable(nil, successorAddr)
          //backupStoreSuc = msg.Store
          // fmt.Println("Found new successor with address: ", successorAddr)
          // send a positive msg back so it knows we accepted proposal and it sets its predecessor
          responseMsg := CommandMessage{"_resProposal", myAddr, msg.SourceAddr, "predecessor", strconv.FormatInt(identifier, 10), dataMap, ""}
          b := getJSONBytes(responseMsg)
          sendMessage(msg.SourceAddr, b)

          // COPY FILES
          msg := CommandMessage{"_copyFiles", myAddr, msg.SourceAddr, "", "iden-here", dataMapSuccessor, "fileBackup"}
          b = getJSONBytes(msg)
          sendMessage(msg.SourceAddr, b)
        } else if msg.Key == "predecessor" && predecessor == -1 {
          predecessor, _ = strconv.ParseInt(msg.Val, 10, 64)
          predecessorAddr = msg.SourceAddr
          initFingerTable(nil, successorAddr)
          //backupStorePred = msg.Store
          // fmt.Println("Found new predecessor with address: ", predecessorAddr)
          // PROBABLY WONT NEED THIS STEP FOR ONE WAY STABILIZATION
          // send a positive msg back so it knows we accepted proposal and it sets its successor if needed
          responseMsg := CommandMessage{"_resProposal", myAddr, msg.SourceAddr, "successor", strconv.FormatInt(identifier, 10), dataMap, ""}
          b := getJSONBytes(responseMsg)
          sendMessage(msg.SourceAddr, b)
        } else {
          // fmt.Println("Response proposal message discarded")
        }
      case "_heartbeat":
        if msg.SourceAddr == successorAddr {
          // successor is still alive so all good
          // fmt.Println("SUCCESSOR ALIVE!")
          dataMapSuccessor = msg.Store
          successorAliveChannel <- true
        }
        if msg.SourceAddr == predecessorAddr {
          // predecessor is still alive so all good
          // fmt.Println("PREDECESSOR ALIVE!")
          dataMapPredecessor = msg.Store
          predecessorAliveChannel <- true
        }
      case "_alive?":
        // received an alive query - send back message to tell I'm still here
        // ISSUE: if i remove else if (just use if which I should logically), then
        // I get timeouts (WHY?)
        if predecessorAddr == msg.SourceAddr {
          //backupStorePred = msg.Store
          sendAliveMessage(predecessorAddr)
        } else if successorAddr == msg.SourceAddr {
          //backupStoreSuc = msg.Store
          sendAliveMessage(successorAddr)
        } else {
            // fmt.Println("successorAddr: ", successorAddr)
            // fmt.Println("predecessorAddr: ", predecessorAddr)
            // fmt.Println("Received alive query from an unexpected node with addr: ", msg.SourceAddr)
        }
      case "_getInfo":
        fmt.Println("Received get info for type: ", msg.Type, " from: ", msg.SourceAddr)
        provideInfo(msg, nodeAddr)
      case "_getVal":
        v, haveKey := getVal(msg.Key)
        if haveKey {
          // respond with Value
          responseMsg := CommandMessage{"_resVal", nodeAddr, msg.SourceAddr, msg.Key, v, nil, ""}
          resp, err := json.Marshal(responseMsg)
          checkError(err)
          buf = []byte(resp)
          // connect to source of request and send Value
          sendMessage(msg.SourceAddr, buf)
        } else {
          // send to next best node
          sendToNextBestNode(k, msg)
        }
      case "_resInfo":
        fmt.Println("Received _resInfo from: ", msg.SourceAddr)
        if msg.Type == "ftab" {
          ftab[k] = msg.Val
          fmt.Println("Set finger table entry ", msg.Key, " to ", ftab[k])
        } else if msg.Type == "streamServer" {
          fmt.Println("Received address of chordNode for streaming: ", msg.Val)
          streamServerChannel <- msg.Val
        } else if msg.Type == "file" {
          fmt.Println("Found node which should hold file part: ", msg.Val)
          // TODO: Transfer file segment to this node OR
          // Return file transfer rpc address of this node (?)
          fileTransferChannel <- msg.Val
        } else {
          fmt.Println("I ain't got no type. Bad bitches the only thing that I like")
        }
      case "_setVal":
        _, haveKey := getVal(msg.Key)
        if haveKey {
          // change Value
          store[msg.Key] = msg.Val
          responseMsg := CommandMessage{"_resGen", nodeAddr, msg.SourceAddr, "", "Key Updated", nil, ""}
          resp, err := json.Marshal(responseMsg)
          checkError(err)
          buf = []byte(resp)
          // connect to source of request and send Value
          sendMessage(msg.SourceAddr, buf)
        } else {
          // send to next best node
          sendToNextBestNode(k,msg)
        }
      case "_locPred" :
        if msg.SourceAddr == successorAddr {
          sendPredInfo(msg.SourceAddr, nodeAddr)
        } else {
          // send to next best node (?)
          sendToNextBestNode(k, msg)
        }
      case "_resLocPred":
        // key in this case would hold asking node's identifier
        predecessor, _ = strconv.ParseInt(msg.Key, 10, 64)
        predecessorAddr = msg.Val;
        // fmt.Println("Updated predecessor to: ", predecessorAddr)
      case "_resDisc" :
        successor = GetIdentifier(msg.Val)
        successorAddr = msg.Val
        c <- "okay"
        // fmt.Println("Successor updated to address: ", msg.Val)
        fmt.Println("Successor Identifier is: ", successor)
      case "_discover":
        nodeIdentifier := GetIdentifier(msg.SourceAddr)
        if successor == -1 {
          // fmt.Println("No successor in network. Setting now to new node...")
          ftab[nodeIdentifier] = msg.SourceAddr // TODO: PROBLEM
          // notify new node of its successor (current successor)
          responseMsg := CommandMessage {"_resDisc", nodeAddr, msg.SourceAddr, "", nodeAddr, nil, ""}
          resMsg, err := json.Marshal(responseMsg)
          checkError(err)
          buf := []byte(resMsg)
          sendMessage(msg.SourceAddr, buf)
          // update successor to new node
          successor = nodeIdentifier
          successorAddr = msg.SourceAddr
          // update predecessor too
          predecessorAddr = msg.SourceAddr
          break
        }
        if betweenIdens(successor, identifier, nodeIdentifier) {
          // incoming node belongs between this node and its current successor
          // Update current successor's pred to new node
          sendPredInfo(successorAddr, msg.SourceAddr)
          // Update new node's pred to me (do we really need this since new node explicitly asks for pred)
          sendPredInfo(msg.SourceAddr, myAddr)
          // fmt.Println("New node fits between me and my successor. Updating finger table...")
          ftab[nodeIdentifier] = msg.SourceAddr
          // notify new node of its successor (current successor)
          responseMsg := CommandMessage {"_resDisc", nodeAddr, msg.SourceAddr, "", successorAddr, nil, ""}
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
          sendToNextBestNode(GetIdentifier(msg.SourceAddr), msg)
          break
        }
      case "_upload": // save file at a node
        //fileIdentifier := GetIdentifier(msg.Key) // the filename is stored in msg.Key
        //nodeToSaveAt := findClosestNode(fileIdentifier)

        in, err := os.Open(msg.Val)
        checkError(err)
        defer in.Close()
        out, err := os.Create("./Downloads/")
        checkError(err)
        defer func() {
          cerr := out.Close()
          if err == nil {
              err = cerr
          }
        }()
        if _, err = io.Copy(out, in); err != nil {
          // do nothing, successful
        }
        err = out.Sync()
        checkError(err)
    }
  }
}

/*
* Attempt to join the system given the ip:port of a running node.
*/
func connectToSystem(nodeAddr string, startAddr string) {
  // fmt.Println("Connecting to peer system...")

  // Figure out where I am in the identifier circle.
  conn, err := net.Dial("udp", startAddr)
  checkError(err)

  defer conn.Close()

  locateSuccessor(conn, nodeAddr)
  locatePredecessor(conn)
  <-c

  initFingerTable(conn, nodeAddr)
  printFingerTable()
}

/*
* Returns a json formatted byte array of the input message
*/
func getJSONBytes(message CommandMessage) []byte {
  resp, err := json.Marshal(message)
  checkError(err)
  return []byte(resp)
}

/*
* Sends message with the whole key value store to node with address addr
*/
func sendKeyMap(addr string) {
  msg := CommandMessage{"_resStoreBackup", myAddr, addr, "", "", dataMap, ""}
  buf := getJSONBytes(msg)
  sendMessage(addr, buf)
}

/*
* Sends a message which requests a node's kv store
*/
func getKeyMap(addr string) {
  msg := CommandMessage{"_storeBackup", myAddr, addr, "", "", nil, ""}
  buf := getJSONBytes(msg)
  sendMessage(addr, buf)
}

/*
* Periodically replicates successor nodes backup to survive loss of data due to node failures
*/
func maintainBackup() {
  for {
    if successorAddr != "" {
      getKeyMap(successorAddr)
    }
    time.Sleep(10 * time.Second)
  }
}

func copyFiles(m map[string]VidFrames) {
  for _, vf := range m {
    for filename, data := range vf.Data {
        path := "FFMPEG/NodesData/" + nodename + "/sample/" + filename
        err := ioutil.WriteFile(path, data, 0644)
        checkError(err)
        fmt.Println("WRITTEN FILES!")
    }
  }
}

// key is filename/foldername and val is the segment sequence number this node holds

func SetStoreVal(filename string) {
  // WILL NEED TO DO THIS LATER TODO
  // iden := GetIdentifier(filename)

  // if iden == identifier {
  //   store[filename] = "available"
  // } else {

  // }
  for store == nil {
    // fmt.Println("Store is null")
  }
  store[filename] = "available"
  // fmt.Printf("Set %s to %s\n", filename, "available")
}

// TODO: return all keys and values this node holds(?)

func GetStoreVal(key string) string {
  return store[key]
}

func GetTransferFileSegmentAddr(filename string) string {
  if successorAddr == "" && predecessorAddr == "" {
    // no one else in the system
    fmt.Println("No one in the system to transfer segment to")
    return ""
  }
  iden := GetIdentifier(filename)
  getNodeInfo(myAddr, iden, "file")
  fmt.Println("Sent node info message for file: ", filename)

  addr := <- fileTransferChannel
  fmt.Println("File transfer chord address received: ", addr)

  msg := CommandMessage{"_fileProposal", myAddr, addr, filename, "iden-here", nil, "file"}
  b := getJSONBytes(msg)
  sendMessage(addr, b)
  addr = <- fileTransferChannel
  fmt.Println("File transfer RPC address received: ", addr)

  return addr
}

func SaveToStore(foldername string, filename string, data []byte) {
  //   Name string
  // FrameStart string
  // TotalFrames int64
  // Data map[string][]byte
  //fmt.Println(string(dataMap[foldername]))
  fmt.Println("Saving to store")
  if reflect.DeepEqual( dataMap[foldername], VidFrames{} ) {
    // no entry currently
    fmt.Println("No entry in data map. Making a new one")
    mp := make(map[string][]byte)
    mp[filename] = data
    vidFrame := VidFrames{foldername, filename, 1, mp} // TODO: not always 1
    dataMap[foldername] = vidFrame
  } else {
    // already an entry for this video - add stuff
    //dataMap[foldername].TotalFrames += 1
    fmt.Println("Entry exists")
    dataMap[foldername].Data[filename] = data
  }
}

func GetStreamingServer(filename string) string {
  if successorAddr == "" && predecessorAddr == "" {
    // no one else in the system
    return streamServerAddress
  }
  //arr := strings.Split(filename, " ")
  iden := GetIdentifier(filename)
  getNodeInfo(myAddr, iden, "streamServer")
  fmt.Println("Sent node info message")
  addr := <- streamServerChannel
  fmt.Println("Address of chord node which will stream: ", addr)

  // now ask the node to prepare stream for this node
  msg := CommandMessage{"_stream", myAddr, addr, "", streamClientAddress, nil, "streamServer"}
  b := getJSONBytes(msg)
  sendMessage(addr, b)

  addr = <- streamServerChannel
  fmt.Println("Address of streaming server: ", addr)
  return addr
}

func Start(thisAddr string, startNodeAddr string, ssa string, sca string, ftAddr string, name string) {
  // Handle the command line.
  //if len(os.Args) != 3 {
  // fmt.Println("Usage: go run node.go [node ip:port] [starter-node ip:port]")
  //  os.Exit(-1)
  //} else {
    myAddr = thisAddr // ip:port of this node
    startAddr := startNodeAddr // ip:port of initial node
    streamServerAddress = ssa
    streamClientAddress = sca
    fileTransferAddr = ftAddr
    nodename = name

    store = make(map[string]string)
    ftab = make(map[int64]string)
    dataMap = make(map[string]VidFrames)
    m = 3
    replicationFactor = 1
    successor = -1
    successorAddr = ""
    predecessor = -1
    predecessorAddr = ""

    c = make(chan string)
    identifier = GetIdentifier(myAddr)

    successorAliveChannel = make(chan bool, 1)
    predecessorAliveChannel = make(chan bool, 1)
    streamServerChannel = make(chan string, 1)
    fileTransferChannel = make(chan string, 1)

    // fmt.Println("THIS NODE'S IDENTIFIER IS: ", identifier)

    if (myAddr == startAddr) {
      // fmt.Println("First node in system. Listening for incoming connections...")
      go startUpSystem(myAddr)
    } else {
        go startUpSystem(myAddr)
        go connectToSystem(myAddr, startAddr)
    }
  //}
  for {
    runtime.Gosched()
  }
}