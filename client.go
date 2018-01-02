package danmu

import (
	"github.com/gorilla/websocket"
	"sync"
)

type (
	cid int64
)

var (
	clientBucket *ClientBucket
)

type Client struct {
	ClientId cid //remain
	Conn     *websocket.Conn
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

type ClientBucket struct {
	Clients map[*websocket.Conn]*Client
	lck     sync.RWMutex
	Room    *Room
}

func InitClientBucket() {
	clientBucket = &ClientBucket{
		Clients: make(map[*websocket.Conn]*Client),
		lck:     sync.RWMutex{},
	}
}

func (cb *ClientBucket) Get(conn *websocket.Conn) (*Client, error) {
	if v, ok := cb.Clients[conn]; ok {
		return v, nil
	} else {
		return nil, ErrRoomDoesNotExist
	}

}

func (cb *ClientBucket) Add(client *Client) error {
	cb.lck.Lock()
	cb.Clients[client.Conn] = client
	cb.lck.Unlock()
	return OK
}

func (cb *ClientBucket) Remove(conn *websocket.Conn) error {
	if _, ok := cb.Clients[conn]; !ok {
		return ErrRoomDoesNotExist
	}
	delete(cb.Clients, conn)
	return OK
}
