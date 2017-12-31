package danmu

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//TODO kafka
const ROOMID_FILED = "room"

var
(
	clients   = make(map[*Client]bool) // connected clients
	broadcast = make(chan Proto)       // broadcast channel
	rwLock    sync.RWMutex             // 读写锁
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}, // 解决域不一致的问题
	}                                  // 将http升级为websocket
)

func onConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Conn: conn}
	roomId := r.FormValue(ROOMID_FILED)
	if roomId == "" {
		client.WriteErrorMsg("incorrect roomId.")
		client.Close()
		return
	}
	roomIdi, err := strconv.Atoi(roomId)
	if err != nil {
		fmt.Println(1)
		client.WriteErrorMsg("incorrect roomId.")
		client.Close()
		log.Println(err)
		return
	}
	room, err := roomBucket.Get(rid(roomIdi))
	if err != nil {
		fmt.Println(2)
		client.WriteErrorMsg("incorrect roomId.")
		client.Close()
		log.Println(err)
		return
	}
	room.AddClient(client)

	//send
	go onMessage(client)

	kali, err := strconv.Atoi(GetConfig("sys", "keepalive_timeout"))
	if err != nil {
		log.Println(err)
		return
	}
	if kali > 0 {
		keepAlive(client, time.Duration(kali)*time.Second)
	}
}

func onMessage(client *Client) {
	defer onClose(client)
	for {
		msgs := Proto{}
		err := client.Conn.ReadJSON(&msgs)
		if err != nil {
			log.Println(err)
			//if check := websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived); check {
			return
			//}
		}
		fmt.Println(msgs)
		broadcast <- msgs
	}
}

func keepAlive(c *Client, timeout time.Duration) {
	lastResponse := time.Now()
	c.Conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()

		return nil
	})
	go func() {
		for {
			err := c.Conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				log.Println("Ping pong timeout, close client", c)
				onClose(c)
				return
			}
		}
	}()
}

func messagePusher() {
	for {
		msg := <-broadcast
		roomId := msg.RoomId
		room, err := roomBucket.Get(rid(roomId))
		if err != nil {
			log.Println(err)
			continue
		}
		for client, _ := range room.GetClients() {
			client.Write(msg)
		}
	}
}

//TODO 支持CloseHandler
func onClose(client *Client) {
	client.Close()
	rwLock.Lock()
	delete(clients, client)
	rwLock.Unlock()
}

func StartServer() {
	InitRoomBucket()

	// http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/ws", onConnect)

	go messagePusher()

	addr := ":" + GetConfig("sys", "port")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
