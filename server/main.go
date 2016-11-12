package main

import (
	"sync"
	"net"
	"log"
	"io"
	"time"
	"bytes"
)

var (
	ServerIP = "192.168.1.235:6263"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)

	//recieve new clipboard item
	go func() {

		ln, err := net.Listen("tcp", "localhost:6263")
		if err != nil {
			log.Fatal(err)
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}

			//blocking read
			var buf bytes.Buffer
			io.Copy(&buf, conn)

			if buf.Len() == 0 {continue} else {

				log.Println("Setting Clipboard")
				//Send new item to all clients (just 1 fixed 1 one for now)
				outconn, err := net.Dial("tcp", "192.168.1.235:6263")
				if err != nil {
					log.Println(err)
					continue
				}

				outconn.Write(buf.Bytes())

			}
			time.Sleep(1 * time.Second)
		}

		wg.Done()
	}()

	wg.Wait()
}
