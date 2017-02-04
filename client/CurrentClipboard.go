package client

import "sync"

//CurrentClipboard holds the current clipboard information
type CurrentClipboard struct {
	text  string
	mutex sync.Mutex
}

//SetText applies the Text string to the systems clipboard
func (cb *CurrentClipboard) SetText(Text string) {
	cb.mutex.Lock()

	cb.text = Text

	cb.mutex.Unlock()
}

//GetText reads the systems clipboard and applies it to the CurrentClipboard
func (cb *CurrentClipboard) GetText() string {
	cb.mutex.Lock()

	Text := cb.text

	cb.mutex.Unlock()

	return Text
}
