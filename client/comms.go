package client

import (
	"log"
	"net"
	"time"

	"errors"

	"github.com/atotto/clipboard"
)

var (
	clientList []string
	cb         = CurrentClipboard{}
)

const (
	maxReadBuffer = 8192
)

//Run looks and listens for cloud-clipboard clients
func Run() {

	//debug purposes only
	Conf.Delete()
	//

	Conf.Load()

	go func() {
		//get remote clipboard
		err := receiveClipboard()
		if err != nil {
			log.Println(err)
		}
	}()

	//get multicasts
	go func() {
		for {
			err := listenForClients()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	//send multicasts
	for {
		err := lookForClients()
		if err != nil {
			log.Println(err)
		}
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

	log.Println("Looking for clients.")
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

	log.Println("Looking for servers.")
	for {
		//log.Println("LFC - Sending Ping")
		_, err := c.Write([]byte("add"))
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

	body := string(b[:n])
	if body == "remove" {
		log.Println("Removing client", src)
		StringArrayRemove(&clientList, src.IP.String())
	} else if body == "add" {
		if !StringArrayContains(clientList, src.IP.String()) { //add if it isnt already in
			clientList = append(clientList, src.IP.String())
			//log.Println("Adding client", src.IP.String())
			go handleClient(src.IP.String())
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

	log.Println("Listening to :6264")
	log.Println("Waiting for a connection...")

	//listen for connections
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer Close(conn)

	log.Println("got a connection!")

	//serve the connection
	for {
		time.Sleep(1 * time.Second)

		buffer := make([]byte, 20000)
		buffSlice := []byte{}

		//read from the connection
		_, err := conn.Read(buffer)
		if err != nil {
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

				log.Println("Setting Clipboard to", msg)

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

func serveClipboard(serverIP string) error {
	conn, err := net.Dial("tcp", serverIP+":6263")
	if err != nil {
		return err
	}

	for {
		time.Sleep(1 * time.Second)

		ReadClipBoard, err := clipboard.ReadAll()
		if err != nil {
			log.Println(err)
			continue
		}

		if ReadClipBoard != cb.GetText() {

			log.Println("Sending Clipboard")
			cb.SetText(ReadClipBoard)

			_, err := conn.Write(AddAuthTCP(ReadClipBoard))
			if err != nil {
				return err
			}
		}
	}
}

func handleClient(serverIP string) {

	//send clipboard to clients
	for {
		time.Sleep(1 * time.Second)

		err := serveClipboard(serverIP)
		if err != nil {
			//check if the error means that the target is offline
			//log.Println(err, reflect.TypeOf(err))
			log.Println(err)
			//stop handling this client for the below reasons
			if err.Error() == "not authed" { //client is not allowed to change the local clipboard
				return
			} //else if err == client is offline {}
		}
	}
}
