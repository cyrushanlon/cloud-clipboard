package client

import "sync"

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