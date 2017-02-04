package client

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/atotto/clipboard"
)

var (
	clientList []string
	cb         = CurrentClipboard{}
)

const (
	address       = "232.49.101.200:6964"
	maxReadBuffer = 8192
)

//ListenForClients waits for a UDP packet to come in and registers/removes the client as required
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

//LookForClients sends out UDP packets in the hope of finding other clients
func LookForClients() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)

	log.Println("Looking for servers.")
	for {
		//log.Println("LFC - Sending Ping")
		c.Write([]byte("add"))
		time.Sleep(1 * time.Second)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	if src.IP.String() == GetLocalIP() {
		return
	}
	//log.Println(n, "bytes read from", src)
	//log.Println()
	body := string(b[:n])
	if body == "remove" {
		log.Println("Removing client", src)
		StringArrayRemove(&clientList, src.String())
	} else if body == "add" {
		if !StringArrayContains(clientList, src.String()) { //add if it isnt already in
			clientList = append(clientList, src.String())
			log.Println("Adding client", src)
			go handleClient(src.IP.String())
		}
	}
}

func receiveClipboard(serverIP string) {
	ln, err := net.Listen("tcp4", ":6264")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening to :6264")

	for {
		log.Println("Waiting for a connection...")

		conn, err := ln.Accept()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			continue
		}

		log.Println("got a connection!")

		for {
			//blocking read
			buffer := make([]byte, 20000)
			buffSlice := []byte{}
			_, err := conn.Read(buffer)

			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}

			for k, v := range buffer {
				if v == 0 {
					buffSlice = buffer[:k]
					break
				}
			}

			if len(buffSlice) == 0 {
				break
			} else {
				log.Println("Setting Clipboard to", string(buffSlice))

				clipboard.WriteAll(string(buffSlice))
				cb.SetText(string(buffSlice))

				time.Sleep(1 * time.Second)
			}

		}
	}
}

func serveClipboard(serverIP string) {
	for {
		conn, err := net.Dial("tcp", serverIP+":6263")
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			ReadClipBoard, err := clipboard.ReadAll()
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Second)
				continue
			}

			if ReadClipBoard != cb.GetText() {

				log.Println("Sending Clipboard")
				cb.SetText(ReadClipBoard)
				conn.Write([]byte(ReadClipBoard))
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func handleClient(serverIP string) {
	go receiveClipboard(serverIP)
	serveClipboard(serverIP)
}
