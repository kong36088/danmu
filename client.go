package danmu

import (
	log "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"time"
)

type (
	cid int64
)

var (
	clientBucket *ClientBucket
	keepalive    int
)

type Client struct {
	ClientId cid //remain
	RoomId   rid
	Conn     *websocket.Conn
}

func NewClient(roomId rid, conn *websocket.Conn) *Client {
	c := new(Client)
	c.RoomId = roomId
	c.Conn = conn
	return c
}

func (client *Client) ReadJSON(v interface{}) error {
	err := client.Conn.ReadJSON(v)
	return err
}

func (client *Client) Read(v interface{}) (int, []byte, error) {
	msgType, msg, err := client.Conn.ReadMessage()
	return msgType, msg, err
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

func (client *Client) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return client.Conn.WriteControl(messageType, data, deadline)
}

func (client *Client) Write(proto *Proto) error {
	err := client.Conn.WriteJSON(proto)
	return err
}

func (client *Client) Close() error {
	return client.Conn.Close()
}

func (client *Client) ErrorReport(err error, msg string) error {
	if msg != "" {
		return client.WriteErrorMsg(msg)
	}
	return OK
}

//keepAlive send ping message to client
// should be called by goroutines
func (client *Client) keepAlive() {
	if keepalive <= 0 {
		return
	}
	timeout := time.Duration(keepalive) * time.Second

	lastResponse := time.Now()
	client.Conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()

		return nil
	})
	go func() {
		for {
			err := client.WriteControl(websocket.PingMessage, []byte("keepalive"), time.Now().Add(writeWait))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				log.Info("Ping pong timeout, close client", client)
				cleaner.CleanClient(client)
				return
			}
		}
	}()
}
func (client *Client) CloseHandler() func(code int, text string) {
	return func(code int, text string) {
		message := []byte{}
		if code != websocket.CloseNoStatusReceived {
			message = websocket.FormatCloseMessage(code, "")
		}
		client.WriteControl(websocket.CloseMessage, message, time.Now().Add(writeWait))
		cleaner.CleanClient(client)
	}
}

type ClientBucket struct {
	Clients map[*websocket.Conn]*Client
	lck     sync.RWMutex
}

func InitClientBucket() error {
	var (
		err error
	)
	clientBucket = &ClientBucket{
		Clients: make(map[*websocket.Conn]*Client),
		lck:     sync.RWMutex{},
	}

	keepalive, err = strconv.Atoi(Conf.GetConfig("sys", "keepalive_timeout"))
	if err != nil {
		return err
	}
	return OK
}

func (cb *ClientBucket) Get(conn *websocket.Conn) (*Client, error) {
	cb.lck.RLock()
	defer cb.lck.RUnlock()
	if v, ok := cb.Clients[conn]; ok {
		return v, OK
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
	cb.lck.Lock()
	delete(cb.Clients, conn)
	cb.lck.Unlock()
	return OK
}
