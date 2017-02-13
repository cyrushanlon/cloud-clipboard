package client

import (
	"errors"
	"net"
	"time"

	"strings"

	"github.com/atotto/clipboard"
)

var (
	clientList map[string]*Client
	cb         = CurrentClipboard{}
)

const (
	maxReadBuffer = 8192
)

//Run looks and listens for cloud-clipboard clients
func Run() {

	//debug purposes only
	//Conf.Delete()
	//log.SetLevel(log.DebugLevel)
	//

	clientList = make(map[string]*Client)
	Conf.Load()

	go func() {
		for {
			//get remote clipboard
			err := receiveClipboard()
			if err != nil {
				//stop handling client if we recieve a connection reset error
				LogWarn(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//get multicasts
	go func() {
		for {
			err := listenForClients()
			if err != nil {
				LogWarn(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//send multicasts
	for {
		err := lookForClients()
		if err != nil {
			LogWarn(err)
		}
		time.Sleep(1 * time.Second)
	}
}

//listenForClients waits for a UDP packet to come in and registers/removes the client as required
func listenForClients() error {
	addr, err := net.ResolveUDPAddr("udp", Conf.MulticastIP)
	if err != nil {
		return err //log.Println(err)
	}
	c, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err //log.Println(err)
	}
	err = c.SetReadBuffer(maxReadBuffer)
	if err != nil {
		return err //log.Println(err)
	}

	LogInfo("Looking for clients.")
	for {
		b := make([]byte, maxReadBuffer)
		n, src, err := c.ReadFromUDP(b)
		if err != nil {
			return err //log.Fatal("ReadFromUDP failed:", err)
		}
		msgHandler(src, n, b)
	}
}

//lookForClients sends out UDP packets in the hope of finding other clients
func lookForClients() error {
	addr, err := net.ResolveUDPAddr("udp", Conf.MulticastIP)
	if err != nil {
		return err //log.Println(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err //log.Println(err)
	}

	LogInfo("Looking for servers.")
	for {
		//log.Println("LFC - Sending Ping")
		_, err := c.Write(AddAuthTCP("add"))
		if err != nil {
			return err //log.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	if src.IP.String() == GetLocalIP() {
		return
	}

	if authed, body := IsAuthedUDP(string(b[:n])); authed {
		if body == "remove" {
			LogInfo("Removing client", src)

			//TODO:
			//review if this should exist
			//Expand this to also terminate a currently handled connection

			//StringArrayRemove(&clientList, src.IP.String())
		} else if body == "add" {

			if _, ok := clientList[src.IP.String()]; !ok {
				clientList[src.IP.String()] = &Client{
					MsgChan: make(chan string),
					IP:      src.IP.String(),
				}

				go clientList[src.IP.String()].Handle()
			}
			/*
				if !StringArrayContains(clientList, src.IP.String()) { //add if it isnt already in
					clientList = append(clientList, src.IP.String())
					//log.Println("Adding client", src.IP.String())
					go handleClient(src.IP.String())
				}
			*/
		}
	}
}

func receiveClipboard() error {
	//listen for packets
	ln, err := net.Listen("tcp4", ":6263")
	if err != nil {
		return err
	}
	defer Close(ln)

	LogInfo("Listening to :6264")
	LogInfo("Waiting for a connection...")

	//listen for connections
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer Close(conn)

	remoteIP := strings.Split(conn.RemoteAddr().String(), ":")[0]

	LogInfo("Connected to", remoteIP)

	//serve the connection
	for {
		time.Sleep(1 * time.Second)

		buffer := make([]byte, 20000)
		buffSlice := []byte{}

		//read from the connection
		_, err := conn.Read(buffer)
		if err != nil {
			//the connection is most likely dead
			switch err.(type) {
			case *net.OpError:
				if v, ok := clientList[remoteIP]; ok {
					v.MsgChan <- "close"
				}
			}
			return err
		}

		//convert the array to a slice
		for k, v := range buffer {
			if v == 0 {
				buffSlice = buffer[:k]
				break
			}
		}

		if len(buffSlice) == 0 {
			//break
		} else {

			msgRaw := string(buffSlice)

			//check that the clipboard change is from an authorised client
			if authed, msg := IsAuthedTCP(msgRaw); authed {

				LogInfo("Setting Clipboard to", msg)

				err := clipboard.WriteAll(msg)
				if err != nil {
					return err
				}
				cb.SetText(msg)

			} else {
				return errors.New("not authed")
			}
		}
	}
}
