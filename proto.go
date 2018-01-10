package danmu

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
)

const (
	OPErr = "error"
	OPMsg = "message"
	OPCountQuery = "count_query"
)

type Proto struct {
	OP    string `json:"op"`
	Message string `json:"message"`
	RoomId  rid    `json:"room_id"`
}

func NewProto() *Proto{
	return &Proto{}
}

func (p Proto) String() string {
	return p.JsonEncode()
}

func (p *Proto) JsonEncode() string {
	j, err := json.Marshal(p)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(j)
}
