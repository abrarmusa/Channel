package main

import "./streamerClient"

func main() {
	streamerClient.ListenForStream("udp://127.0.0.1:1234")
	handler0 := streamerClient.GetRpcHandler(":1342")
	handler1 := streamerClient.GetRpcHandler(":1354")
	streamerClient.StartStreaming(handler0, 0, "0")
	streamerClient.StartStreaming(handler1, 1, "300")
}