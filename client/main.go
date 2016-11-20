package main

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/atotto/clipboard"

	"bytes"
	"time"
)

var (
	ServerIP = "localhost:6264"
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

		ln, err := net.Listen("tcp", "localhost:6263")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Listening to localhost:6263")

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}

			log.Println("got a connection")

			//blocking read
			var buf bytes.Buffer
			io.Copy(&buf, conn)

			if buf.Len() == 0 {
				continue
			} else {

				log.Println("Setting Clipboard")

				clipboard.WriteAll(buf.String())
				cb.SetText(buf.String())
			}
			time.Sleep(1 * time.Second)
		}

		wg.Done()
	}()

	//serve clipboard out
	go func() {

		conn, err := net.Dial("tcp", ServerIP)

		if err != nil {
			log.Fatal(err)
		}

		for {
			ReadClipBoard, err := clipboard.ReadAll()
			if err != nil {
				log.Println(err)
				continue
			}

			if ReadClipBoard != cb.GetText() {

				log.Println("Sending Clipboard")
				cb.SetText(ReadClipBoard)
				conn.Write([]byte(ReadClipBoard))
			}

			time.Sleep(1 * time.Second)
		}

		wg.Done()
	}()

	wg.Wait()

}
