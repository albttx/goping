package goping

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/icmp"
)

type Pinger struct {
	hostname string
	hostip   string
	socket   *icmp.PacketConn
}

func newPinger(addr string) *icmp.PacketConn {
	ip, errIp := net.LookupIP(addr)
	if errIp != nil {
		log.Print("Bad ip")
	}

	fmt.Println(ip[0].String())
	sock, err := icmp.ListenPacket("ip4:icmp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return sock
}

func ping() {

}

func listener() {

}

func Ping(addr string) {
	newPinger(addr)
}
