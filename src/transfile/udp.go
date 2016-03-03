package transfile

import (
	"net"
)

func SetupUdpPort() string {
	thisaddr, err := net.ResolveUDPAddr("udp", ":0")
	checkError(err)
	conn, err := net.ListenUDP("udp", thisaddr)
	checkError(err)
	return conn.LocalAddr().String()[4:]
}
