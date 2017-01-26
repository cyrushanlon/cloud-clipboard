package client

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

const (
	address       = "232.49.101.200:6964"
	maxReadBuffer = 8192
)

//holds a list of active clients
var ClientList []string

func StringArrayContains(list []string, target string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == target {
			return true
		}
	}
	return false
}

func StringArrayRemove(list *[]string, target string) {
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == target {
			l := append((*list)[:i], (*list)[i+1:]...)
			list = &l
			return
		}
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	//log.Println(n, "bytes read from", src)
	//log.Println()
	body := hex.Dump(b[:n])
	if body == "remove" {
		StringArrayRemove(&ClientList, src.String())
	} else if body == "add" {
		if !StringArrayContains(ClientList, src.String()) { //add if it isnt already in
			ClientList = append(ClientList, src.String())
		}
	}
}

func ListenForClients() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	c.SetReadBuffer(maxReadBuffer)

	log.Println("Looking for clients.")
	for {
		b := make([]byte, maxReadBuffer)
		n, src, err := c.ReadFromUDP(b)
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

	log.Println("Looking for servers.")
	for {
		log.Println("LFC - Sending Ping")
		c.Write([]byte("add"))
		time.Sleep(1 * time.Second)
	}
}
