package danmu

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func (client *Client) WriteErrorMsg(msg string) error {
	err := client.Conn.WriteJSON(Proto{
		OP:      OPErr,
		Message: msg,
		RoomId:  -1,
	})
	return err
}

func (client *Client) WriteMessage(msg string, roomId rid) error {
	err := client.Conn.WriteJSON(Proto{
		OP:      OPMsg,
		Message: msg,
		RoomId:  roomId,
	})
	return err
}

func (client *Client) Write(proto Proto) error {
	err := client.Conn.WriteJSON(proto)
	return err
}

func (client *Client) Close() {
	client.Conn.Close()
}
