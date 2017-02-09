package client

import "log"

//Closer is used for closing anything that has a Close function
//It can be used to defer the close and logs any errors found
type Closer interface {
	Close() error
}

//Close is used for closing something and logging any errors.
func Close(c Closer) {
	if c != nil {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
