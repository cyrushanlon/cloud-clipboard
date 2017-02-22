package main

import "github.com/cyrushanlon/cloud-clipboard"

func main() {
	//check if the process is already open and if it is show
	//maybe use command line args to do that?

	//debug purposes only
	client.Conf.Delete()
	//log.SetLevel(log.DebugLevel)
	//

	//run the client
	client.Conf.Load()

	go client.RunLocal()
	client.RunRemote()
}
