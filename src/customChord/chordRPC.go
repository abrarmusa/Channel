package main

import (
	"fmt"
	"math"
	"net"
	"net/rpc"
	"os"
	"math/big"
	"crypto/sha1"
  	"encoding/hex"
  	"runtime"
  	"time"
)

type (


	ChordService int 		// type for rpc

	Msg struct {
		SourceAddress 		string 				// source address
		Key 				string
		KeyIdentifier 		int64
		KeyType 			string 				// stores inquired key's type
		Val 				string 				// holds any value that the client wants the server to use
	}

	Reply struct {
		Key 				string
		Val 				string
	}

)

var (
	nodeAddress				string 				// rpc address this node listens on
	peerAddress 			string 				// another node's address to connect to
	nodeIdentifier 			int64
	ftab 					map[int64]string
	successorIdentifier 	int64
	predecessorIdentifier	int64
	successorAddress		string	
	predecessorAddress		string
	m 						float64
	rpcChain				chan string
)

func main () {

	// set this node's rpc address TODO
	nodeAddress = os.Args[1]
	peerAddress = os.Args[2]

	// set m
	m = 3

	// init rpc chain channel
	rpcChain = make(chan string, 1)

	nodeIdentifier = getIdentifier(nodeAddress)
	// initialize finger table
	ftab = make(map[int64]string)

	go launchRPCService()

	if nodeAddress == peerAddress {
		fmt.Printf("First node %s joining the system\n", nodeAddress)
		successorAddress = ""
		predecessorAddress = ""
	} else {
		fmt.Printf("Connecting to peer %s\n", peerAddress)
		// start communicating with starter node
		// find successor
		// find predecessor

		// send GetKeyInfo message for type node to get discovered
		var reply Reply
		msg := Msg {nodeAddress, nodeAddress, nodeIdentifier, "node", ""}
		handler := getRpcHandler(peerAddress)
		err := handler.Call("ChordService.GetKeyInfo", &msg, &reply)
		checkError(err)
		fmt.Printf("Reply received for GetKeyInfo: %s\n", reply.Val)

		// wait to get successor and predecessor
		for successorAddress == "" || predecessorAddress == "" {
			fmt.Println("No successor and predecessor addresses. Waiting ...")
			time.Sleep(2 * time.Second)
		}

		// populate finger table
		populateFingerTable()

		printFingerTable()
	}

	for {
		runtime.Gosched()
	}

	// find successor and predecessor (?)
	// by simply connecting to system - call GetKeyInfo (?)

}

//////////////////////////////////////////////////////
/*				RPC FUNCTIONS (INBOUND)				*/
//////////////////////////////////////////////////////

// func (this *ChordService) Alive(msg *Msg, reply *Reply) {
// 	// TODO
// 	// could we just detect tcp close to see if our successor/predecessor is dead?
// }


	// Msg struct {
	// 	SourceAddress 		string 				// source address
	// 	Key 				string
	// 	KeyIdentifier 		int64
	// 	KeyType 			string 				// stores inquired key's type
	// 	Val 				string 				// holds any value that the client wants the server to use
	// }

func (this *ChordService) GetKeyInfo(msg *Msg, reply *Reply) error {
	fmt.Println("Received GetKeyInfo message: ", msg)
	// check if key's identifier lies between me and my successor
		// if it does then:
			// if the key is a node then it falls between me and my successor (updates required - node join)
			// if the key is a file then reply with successor's address cause it holds the file
		// else forward to next best node in our finger table (closest to key's identifier/max identifer in ftab less than key's identifier)
	if successorAddress == "" && predecessorAddress == "" {
		// only node in system - deal accordingly
		// ask new node to set me as a successor and a predecessor
		//fmt.Println("Found another node. Not lonely anymore")
		var reply Reply
		msg0 := Msg {nodeAddress, "", -1, "", nodeAddress}
		handler := getRpcHandler(msg.SourceAddress)
		err := handler.Call("ChordService.SetPredecessor", &msg0, &reply)
		checkError(err)
		err = handler.Call("ChordService.SetSuccessor", &msg0, &reply)
		checkError(err)
		//fmt.Printf("Reply received for SetPredecessor: %s\n",reply.Val)
		// set new node as my successor and predecessor
		successorAddress = msg.SourceAddress
		successorIdentifier = getIdentifier(msg.SourceAddress)
		predecessorAddress = msg.SourceAddress
		predecessorIdentifier = getIdentifier(msg.SourceAddress)
		ftab[nodeIdentifier+1] = msg.SourceAddress

		populateFingerTable()
		printFingerTable()

	} else if msg.KeyIdentifier == nodeIdentifier {
		// looking for me
		fmt.Println("Someone inquired about my identifier. Sending info back.")
		reply.Val = nodeAddress
	} else if betweenIdentifiers(msg.KeyIdentifier) {
		if msg.KeyType == "node" {
			// ask new node to select me as its predecessor and my old successor as its successor TODO
			// Need: SetPredecessor(), SetSuccessor() - make rpc calls
			fmt.Println("BETWEEN ME AND MY successor")
			var reply Reply
			msg0 := Msg {nodeAddress, "", -1, "", nodeAddress}
			handler := getRpcHandler(msg.SourceAddress)
			err := handler.Call("ChordService.SetPredecessor", &msg0, &reply)
			checkError(err)
			//fmt.Printf("Reply received for SetPredecessor: %s\n",reply.Val)

			msg0 = Msg {nodeAddress, "", -1, "", successorAddress}
			handler = getRpcHandler(msg.SourceAddress)
			err = handler.Call("ChordService.SetSuccessor", &msg0, &reply)
			checkError(err)
			//fmt.Printf("Reply received for SetSuccessor: %s\n",reply.Val)

			// ask my old successor to select new node as its predecessor TODO
			// Need: SetPredecessor() - make rpc call
			msg0 = Msg {nodeAddress, "", -1, "", msg.SourceAddress}
			handler = getRpcHandler(successorAddress)
			err = handler.Call("ChordService.SetPredecessor", &msg0, &reply)
			checkError(err)
			//fmt.Printf("Reply received for SetPredecessor: %s\n",reply.Val)			

			// change my successor and update finger table entry
			successorAddress = msg.SourceAddress
			successorIdentifier = getIdentifier(msg.SourceAddress)
			ftab[nodeIdentifier+1] = msg.SourceAddress
			reply.Val = "Accepted in the family"
		} else {
			// file or ftab population inquiry - simply send successor's address
			fmt.Println("Between me and my successor: File or ftab inquiry received")
			fmt.Printf("successorIden: %d\npredecessorIden: %d\n", successorIdentifier, predecessorIdentifier)
			//fmt.Println("Message: ", msg)
			reply.Val = successorAddress
		}
	} else {
		// send to next best node
		sendToNextBestNode(msg)
		reply.Val = <- rpcChain
	}
	return nil

}

func (this *ChordService) SetPredecessor(msg *Msg, reply *Reply) error {
	fmt.Println("Updating predecessor to: ", msg.Val)
	predecessorAddress = msg.Val
	predecessorIdentifier = getIdentifier(msg.Val)
	reply.Val = "Okay"
	//populateFingerTable()
	//printFingerTable()
	return nil
}

func (this *ChordService) SetSuccessor(msg *Msg, reply *Reply) error {
	fmt.Println("Updating successor to: ", msg.Val)
	successorAddress = msg.Val
	successorIdentifier = getIdentifier(msg.Val)
	// adjust finger table
	ftab[nodeIdentifier+1] = msg.Val
	//populateFingerTable()
	printFingerTable()
	return nil
}



//////////////////////////////////////////////////////
/*				RPC FUNCTIONS (INBOUND) END			*/
//////////////////////////////////////////////////////

/* 
* Set up the listener for RPC requests, serve the connections when required.
*/
func launchRPCService() {
  // Set up RPC service
  server := new(ChordService)
  rpc.Register(server)
  rpcAddr, err := net.ResolveTCPAddr("tcp", nodeAddress)
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

/*
* Initializes finger table populating entries from iden+2^0 to iden+2^m
*/
func populateFingerTable() {
	for i := 0; i < int(m); i++ {
		key := int64( math.Mod( float64(nodeIdentifier) + math.Pow(2, float64(i)), math.Pow(2, float64(m)) ) )

		var reply Reply
		msg := Msg {nodeAddress, "", key, "ftab", ""}
		handler := getRpcHandler(successorAddress)
		err := handler.Call("ChordService.GetKeyInfo", &msg, &reply)
		checkError(err)
		fmt.Printf("Reply received for entry %d in ftab: %s\n", key, reply.Val)

		ftab[key] = reply.Val
	}
}

// func initFingerTable(conn net.Conn, nodeAddr string) {
//   thisIden := GetIdentifier(nodeAddr)
//   for i := 0; i < int(m); i++ {
//     key := int64( math.Mod( float64(thisIden) + math.Pow(2, float64(i)), math.Pow(2, float64(m)) ) )
//     getNodeInfo(nodeAddr, key, "ftab")
//   }
// }

func getRpcHandler(rpcAddr string) (*rpc.Client) {
	//var err error
	var nodeRPCHandler *rpc.Client
	//fmt.Println("Dialing address: ", rpcAddr)
	nodeRPCHandler, err := rpc.Dial("tcp", rpcAddr)
	
	checkError(err)
	//defer nodeRPCHandler.Close()
	return nodeRPCHandler
}

/*
* Sends to next best candidate for finding KeyIdentifier by searching through finger table
*/
func sendToNextBestNode(msg *Msg) {
	var closestNode string
	minDistanceSoFar := int64(math.MaxInt64)
	for nodeIden, nodeAddr := range ftab {
	  diff := nodeIden - msg.KeyIdentifier
	  if diff < minDistanceSoFar {
	    minDistanceSoFar = diff
	    closestNode = nodeAddr
	  }
	}
	// send message to closestNode - rpc outbound call TODO
	// (?) Always GetKeyInfo service (?)
  	var reply Reply
	handler := getRpcHandler(closestNode)
	err := handler.Call("ChordService.GetKeyInfo", &msg, &reply)
	checkError(err)
	fmt.Printf("Reply received for rpcChain: %s\n", reply.Val)
	rpcChain <- reply.Val
}

/*
* Returns an int64 identifier for an input key
*/
func getIdentifier(key string) int64 {
  id := computeSHA1Hash(key)
  k := big.NewInt(0)
  if _, ok := k.SetString(id, 16); ok {
    //fmt.Println("Number: ", k)
  } else {
    //fmt.Println("Unable to parse into big int")
  }
  power := int64(math.Pow(2, m))
  ret := (k.Mod(k, big.NewInt(power))).Int64()

  fmt.Println("Identifier for ", key, " : ", ret)

  return ret
}

/*
* Checks if an identifier iden lies between this node and its successor
*/
func betweenIdentifiers(iden int64) bool {
  if successorIdentifier < nodeIdentifier {
    if iden > nodeIdentifier && iden > successorIdentifier {
      return true
    } else if iden < nodeIdentifier && iden < successorIdentifier {
      return true
    }
  } else if successorIdentifier > nodeIdentifier {
    if iden > nodeIdentifier && iden < successorIdentifier {
      return true
    }
  }
  return false
}

/*
* Checks error value and prints/exits if non nil.
*/
func checkError(err error) {
  if err != nil {
    fmt.Println("Error thrown: ", err)
    os.Exit(-1)
  }
}

/*
* Returns the SHA1 hash Value of input key as a string
*/
func computeSHA1Hash(key string) string {
  buf := []byte(key)
  h := sha1.New()
  h.Write(buf)
  str := hex.EncodeToString(h.Sum(nil))
  return str
}

/* Prints the finger table entries to standard output.
 */
func printFingerTable() {
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf(" Finger table (unordered) for this node: %d\n", nodeIdentifier)
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
  fmt.Printf("| ID   |    VAL    |\n")

  // Runs up to size m.
  for id := range ftab {
    fmt.Printf("| %3d  | %9s |\n", id, ftab[id])
  }
  fmt.Println(" -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ ")
}