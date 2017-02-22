package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
)

//RunRemote handles the connections to the remote server
func RunRemote() {
	go setRemoteClipboard()
	getRemoteClipboard()
}

func setRemoteClipboard() {
	for { // dont set remote clipboard if we dont want to
		if !Conf.AllowServer {
			time.Sleep(10 * time.Second)
			continue
		}

		//do it
		time.Sleep(1 * time.Second)

		//res, err := http.Get(remoteURL)
	}
}

func getRemoteClipboard() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{Transport: tr}

	for {
		if !Conf.AllowServer { // dont get remote clipboard if we dont want to
			time.Sleep(10 * time.Second)
			continue
		}

		//do it
		time.Sleep(2 * time.Second)

		//make the request
		req, err := http.NewRequest("GET", Conf.RemoteIP, nil)
		if err != nil {
			LogWarn(err)
			continue
		}
		req.SetBasicAuth(Conf.Username, Conf.Password)

		//do the request
		res, err := cli.Do(req)
		if err != nil {
			LogWarn(err)
			continue
		}

		//was the request succesful?
		if res.StatusCode != 200 {
			continue
		}

		//read the body out
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			LogWarn(err)
			continue
		}

		Close(res.Body)

		//set the local clipboard if required
		newText := string(data)

		if cb.GetText() != newText {

			LogInfo("setting clipbaord to", newText)

			err = clipboard.WriteAll(newText)
			if err != nil {
				LogWarn(err)
				continue
			}
			cb.SetText(newText)
		}
	}
}
