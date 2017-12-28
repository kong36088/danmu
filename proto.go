package danmu

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (msg Message) String() string {
	return msg.JsonEncode()
}

func (msg *Message) JsonEncode() string {
	j, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(j)
}
