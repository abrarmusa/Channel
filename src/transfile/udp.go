package transfile

import (
	"fmt"
	"net"
)

type UdpInfo struct {
	Conn *net.UDPConn
	Port string
}

func (ui *UdpInfo) SetupUdp() {
	addr, err := net.ResolveUDPAddr("udp", ui.Port)
	checkError(err)
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	ui.Conn = conn
	ui.Port = conn.LocalAddr().String()[4:]
}

func (ui *UdpInfo) ReceiveUdpCall() {
	defer ui.Conn.Close()
	buf := make([]byte, 1024)
	for {
		n, _, err := ui.Conn.ReadFromUDP(buf)
		checkError(err)
		fmt.Println(string(buf[0:n]))
	}
}

func (ui *UdpInfo) SendUdpCall(addr_str string) {
	addr, err := net.ResolveUDPAddr("udp", addr_str)
	_, err = ui.Conn.WriteToUDP([]byte("thanks"), addr)
	checkError(err)
}
