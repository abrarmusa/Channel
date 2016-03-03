package transfile

import (
	"net"
)

type UdpInfo struct {
	Conn *net.UDPConn
	Port string
}

func (ui *UdpInfo) SetupUdp() {
	thisaddr, err := net.ResolveUDPAddr("udp", ui.Port)
	checkError(err)
	conn, err := net.ListenUDP("udp", thisaddr)
	checkError(err)
	ui.Conn = conn
	ui.Port = conn.LocalAddr().String()[4:]
}
