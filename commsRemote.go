package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
)

//RunRemote handles the connections to the remote server
func RunRemote() {

	for {
		res, err := http.Get("https://cloud-clipboard-server.herokuapp.com")
		if err != nil {
			log.Println(err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		res.Body.Close()

		log.Println("setting clipbaord to", string(data))

		err = clipboard.WriteAll(string(data))
		if err != nil {
			LogWarn(err)
		}

		cb.SetText(string(data))

		time.Sleep(1 * time.Second)
	}
}
