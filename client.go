package client

import (
	"net"
	"reflect"
	"time"

	"github.com/atotto/clipboard"
)

//Client handles an individual clients connections
type Client struct {
	MsgChan  chan string
	IP       string
	LastSent string
}

//Handle handles the sending of the local clipboard to the client
func (c *Client) Handle() {

	//send clipboard to clients
	for {
		time.Sleep(1 * time.Second)

		err := c.serveClipboard(c.IP)
		if err != nil {
			//check if the error means that the target is offline
			//stop handling this client for the below reasons
			if err.Error() == "not authed" { //client is not allowed to change the local clipboard
				delete(clientList, c.IP)
				return
			}

			switch err.(type) {
			case *net.OpError: //connection is closed at remote end
				LogWarn(err)
				delete(clientList, c.IP)
				return
			}

			LogErr(err, "|", reflect.TypeOf(err))
		} else { // the connection to the client is closed
			delete(clientList, c.IP)
			return
		}
	}
}

func (c *Client) serveClipboard(serverIP string) error {

	conn, err := net.Dial("tcp", serverIP+":6263")
	if err != nil {
		return err
	}
	defer Close(conn)

	for {
		time.Sleep(1 * time.Second)

		select {
		case msg := <-c.MsgChan:
			if msg == "close" {
				return nil
			}
		default:

			if len(clientList) == 0 {
				continue
			}

			ReadClipBoard, err := clipboard.ReadAll()
			if err != nil {
				if err.Error() != "exit status 1" { //clipboard is empty
					return err
				}
				continue
			}

			if ReadClipBoard != c.LastSent {

				LogInfo("Sending Clipboard to", serverIP)

				msg, err := AddAuthTCP(ReadClipBoard)
				if err != nil {
					LogErr(err)
					continue
				}

				_, err = conn.Write(msg)
				if err != nil {
					return err
				}
				//if its ok, set the values
				cb.SetText(ReadClipBoard)
				c.LastSent = ReadClipBoard
			}
		}
	}
}
