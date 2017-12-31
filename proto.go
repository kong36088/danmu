package danmu

import (
	"encoding/json"
	"log"
)

const (
	OPErr = "error"
	OPMsg = "message"
)

type Proto struct {
	OP    string `json:"op"`
	Message string `json:"message"`
	RoomId  rid    `json:"room_id"`
}

func (p Proto) String() string {
	return p.JsonEncode()
}

func (p *Proto) JsonEncode() string {
	j, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(j)
}
