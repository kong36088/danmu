package danmu

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func (client *Client) Close() {
	client.conn.Close()
}
