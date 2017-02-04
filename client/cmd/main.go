package main

import "github.com/cyrushanlon/cloud-clipboard/client"

var (
	serverIP = "192.168.1.236:6263"
)

func main() {
	//send multicasts
	go client.LookForClients()
	//get multicasts
	client.ListenForClients()
}
