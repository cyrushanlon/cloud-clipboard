package client

import (
	"net"
	"reflect"
	"time"
)

type Client struct {
	MsgChan chan string
	IP      string
}

func (c Client) Handle() {

	//send clipboard to clients
	for {
		time.Sleep(1 * time.Second)

		err := serveClipboard(c.IP)
		if err != nil {
			//check if the error means that the target is offline
			//stop handling this client for the below reasons
			if err.Error() == "not authed" { //client is not allowed to change the local clipboard
				return
			}

			switch err.(type) {
			case *net.OpError: //connection is closed at remote end
				LogWarn(err)
				return
			}

			LogErr(err, "|", reflect.TypeOf(err))
		}

	}

}
