package main

import (
	"log"
	"net"
	"sync"
	"time"
)

var (
	ClientIP = "192.168.1.235:6264"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)

	//recieve new clipboard item
	go func() {

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

			log.Println(conn.RemoteAddr().String())

			for {

				buffer := make([]byte, 4096)
				n, err := conn.Read(buffer)
				if err != nil {
					log.Println(err)
					break
				}

				if n == 0 {
					continue
				} else {

					log.Println("Setting Clipboard to", string(buffer))
					//Send new item to all clients (just 1 fixed 1 one for now)
					outconn, err := net.Dial("tcp4", ClientIP)
					if err != nil {
						log.Println(err)
						continue
					}

					n, err := outconn.Write(buffer)
					if err != nil {
						log.Println("n:", n, err)
					}
					outconn.Close()
				}
			}
			time.Sleep(1 * time.Second)
		}

		wg.Done()
	}()

	wg.Wait()
}
