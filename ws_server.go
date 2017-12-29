package danmu

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"sync"
	"time"
	"strconv"
)

//TODO kafka
var
(
	clients   = make(map[*Client]bool) // connected clients
	broadcast = make(chan Message)     // broadcast channel
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
	client := &Client{conn}
	rwLock.Lock()
	clients[client] = true
	rwLock.Unlock()
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
		var msgs Message
		messageType, message, err := client.conn.ReadMessage()
		fmt.Println(message, messageType)
		if err != nil {
			log.Println(err)
			//if check := websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived); check {
			return
			//}
		}
		switch messageType {
		case websocket.TextMessage:
			err = json.Unmarshal([]byte(message), &msgs)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(msgs)
			broadcast <- msgs
		case websocket.CloseMessage:
			return
		default:
			continue
		}

	}
}

func keepAlive(c *Client, timeout time.Duration) {
	lastResponse := time.Now()
	c.conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()

		return nil
	})

	go func() {
		for {
			err := c.conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
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
		rwLock.RLock()
		for client, _ := range clients {
			client.conn.WriteMessage(1, []byte(msg.JsonEncode()))
		}
		rwLock.RUnlock()
	}
}

func onClose(client *Client) {
	client.Close()
	rwLock.Lock()
	delete(clients, client)
	rwLock.Unlock()
}

func StartServer() {
	http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/ws", onConnect)

	go messagePusher()

	addr := ":" + GetConfig("sys", "port")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
