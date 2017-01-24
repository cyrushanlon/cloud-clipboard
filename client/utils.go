package client

import (
	"net"
	"log"
	"encoding/hex"
	"time"
)

const (
	address = "232.49.101.200:6964"
	maxReadBuffer = 8192
)

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}

func ListenForClients() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	listener.SetReadBuffer(maxReadBuffer)
	for {
		b := make([]byte, maxReadBuffer)
		n, src, err := listener.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		msgHandler(src, n, b)
	}
}

func LookForClients() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	for {
		c.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}