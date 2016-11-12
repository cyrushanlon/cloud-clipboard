package cl
/*
import (
	"net/http"
	"strings"
	"github.com/atotto/clipboard"
	"time"
	"errors"
)



func SendClipboard() error{

	str, err := clipboard.ReadAll()
	if err != nil {
		return err
	}

	client := http.Client{Timeout:time.Second * 5}

	req, err := http.NewRequest("POST", ServerIP, strings.NewReader(str))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if !resp.StatusCode == 200 {
		return errors.New("Status code not OK")
	}

	return nil
}

func GetClipboard() {



}*/