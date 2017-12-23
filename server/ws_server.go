package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"sync"
)

//TODO 心跳检测 断开检测
//TODO kafka
var clients = make(map[*Client]bool) // connected clients
var broadcast = make(chan Message)   // broadcast channel
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, // 解决域不一致的问题
}                                    // 将http升级为websocket
var rwLock sync.RWMutex              // 读写锁

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
	go listenMessage(client)
}

func listenMessage(client *Client) {
	defer onClose(client)
	for {
		var msgs Message
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			client.conn.Close()
			rwLock.Lock()
			delete(clients, client)
			rwLock.Unlock()
			return
		}
		err = json.Unmarshal([]byte(message), &msgs)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(msgs)
		broadcast <- msgs
	}
}

func messagePusher() {
	for {
		msg := <-broadcast
		rwLock.RLock()
		for client, _ := range clients {
			//TODO messageType
			client.conn.WriteMessage(1, []byte(msg.JsonEncode()))
		}
		rwLock.RUnlock()
	}
}

func onClose(client *Client){
	client.conn.Close()
}

func StartServer() {
	http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/ws", onConnect)

	go messagePusher()

	err := http.ListenAndServe(":9500", nil)
	if err != nil {
		log.Fatal(err)
	}
}
