package goping

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ProtocolICMPv4 = 1
	ProtocolICMPv6 = 48 // defined in internal/iana
)

type Pinger struct {
	hostname string
	hostip   string
	socket   *icmp.PacketConn
	sendTime time.Time
}

func newPinger(addr string) *Pinger {
	var p Pinger

	ip, err := net.LookupIP(addr)
	if err != nil {
		log.Fatal(err)
	}
	p.hostip = ip[0].String()
	p.hostname = addr

	p.socket, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	return &p
}

func (p *Pinger) sendICMP() {
	packet := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
	}
	body := &icmp.Echo{
		ID:   os.Getpid() & 0xffff,
		Seq:  0,
		Data: []byte("toc toc toc"),
	}

	for {
		packet.Body = body
		pkt, err := packet.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}

		p.sendTime = time.Now()
		cc, err := (p.socket).WriteTo(pkt, &net.IPAddr{IP: net.ParseIP(p.hostip)})
		if err != nil || cc == 0 {
			log.Fatal(err)
		}
		body.Seq += 1
		time.Sleep(time.Second * 1)
	}
}

func printEchoReply(icmp_h *icmp.Echo) {
	ip_h, err := ipv4.ParseHeader(icmp_h.Data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" icmp_seq=%-3v ttl=%v", icmp_h.Seq, ip_h.TTL)
}

func (p *Pinger) recvICMP() {
	var recvTime time.Time
	buff := make([]byte, 512)

	for {
		cc, peer, err := (p.socket).ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}

		msg, err := icmp.ParseMessage(ProtocolICMPv4, buff)
		if err != nil {
			log.Fatal(err)
		}
		recvTime = time.Now()
		fmt.Printf("%-3v byes from %s: %.5vms (%v)", cc, peer, recvTime.Sub(p.sendTime), msg.Type)
		if msg.Type == ipv4.ICMPTypeEchoReply {
			printEchoReply(msg.Body.(*icmp.Echo))
		}
		fmt.Printf("\n")
	}
}

func Ping(addr string) {
	p := newPinger(addr)
	go p.sendICMP()
	p.recvICMP()
}
