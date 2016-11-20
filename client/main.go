package main

import (
	"github.com/atotto/clipboard"
	"log"
	"net"
	"sync"
	"time"
	"io"
)

var (
	ServerIP = "192.168.1.237:6263"
)

type CurrentClipboard struct {
	text  string
	mutex sync.Mutex
}

func (cb *CurrentClipboard) SetText(Text string) {
	cb.mutex.Lock()

	cb.text = Text

	cb.mutex.Unlock()
}

func (cb *CurrentClipboard) GetText() string {
	cb.mutex.Lock()

	Text := cb.text

	cb.mutex.Unlock()

	return Text
}

func main() {

	var wg sync.WaitGroup

	cb := CurrentClipboard{}

	wg.Add(2)

	//serve clipboard in
	go func() {

		ln, err := net.Listen("tcp4", ":6264")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Listening to :6264")

		for {
			conn, err := ln.Accept()
			if err != nil {
				if err != io.EOF {
					log.Println( err)
				}
				continue
			}

			log.Println("got a connection")

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
					log.Println("Connection.Read got nothing")
					break
				} else {
					log.Println("Setting Clipboard to", string(buffSlice), len(buffSlice))

					clipboard.WriteAll(string(buffSlice))
					cb.SetText(string(buffSlice))

					time.Sleep(1 * time.Second)
				}

			}
		}

		wg.Done()
	}()

	//serve clipboard out
	go func() {

		for {

			conn, err := net.Dial("tcp", ServerIP)
			if err != nil {
				log.Println(err)
				time.Sleep(1*time.Second)
				continue
			}

			for {
				ReadClipBoard, err := clipboard.ReadAll()
				if err != nil {
					log.Println(err)
					time.Sleep(1*time.Second)
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

		wg.Done()
	}()

	wg.Wait()

}
