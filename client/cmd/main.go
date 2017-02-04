package main

import "github.com/cyrushanlon/cloud-clipboard/client"

func main() {
	//get multicasts
	go client.ListenForClients()
	//send multicasts
	client.LookForClients()
}
