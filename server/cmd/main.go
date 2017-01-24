package main

import (
	"log"
	"net"
	"sync"
	"strings"
)

var (
//	ClientIPs = []string{}
	wg sync.WaitGroup
)

func main() {

	for {

		ln, err := net.Listen("tcp4", ":6263")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Listening to localhost:6263")

		for {
			log.Println("Waiting to accept a connection")
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}

			newAdd := strings.Split(conn.RemoteAddr().String(), ":")[0]

			wg.Add(1)
			log.Println("serving", newAdd)
			go ServeConnection(conn, newAdd)
		}
	}
}

func ServeConnection(conn net.Conn, addr string) {
	//serve the connection
	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			wg.Done()
			break
		}

		if n == 0 {
			continue
		} else {

			log.Println("Setting Clipboard on", addr)

			outconn, err := net.Dial("tcp4", addr + ":6264")
			if err != nil {
				log.Println(err)
				return
			}

			n, err := outconn.Write(buffer)
			if err != nil {
				log.Println("n:", n, err)
			}
			outconn.Close()
/*
			var wg2 sync.WaitGroup

			wg2.Add(len(ClientIPs))
			for _, v := range(ClientIPs) {
				go func() {
					outconn, err := net.Dial("tcp4", v+":6264")
					if err != nil {
						log.Println(err)
						wg2.Done()
						return
					}

					n, err := outconn.Write(buffer)
					if err != nil {
						log.Println("n:", n, err)
					}
					outconn.Close()
					wg2.Done()
				}()
			}
			wg2.Wait()
*/
		}
	}
	wg.Done()
}