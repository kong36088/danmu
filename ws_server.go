package danmu

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	writeWait   = time.Second
	roomIdFiled = "room"
)

//TODO kafka
//TODO log4go
var
(
	broadcast = make(chan Proto) // broadcast channel
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}, // 解决域不一致的问题
	}                            // 将http升级为websocket
)

func onConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(0, conn)
	roomId := r.FormValue(roomIdFiled)
	roomIdi, err := strconv.Atoi(roomId)
	if err != nil {
		client.WriteErrorMsg("incorrect roomId.")
		cleaner.CleanClient(client)
		log.Printf("Parse roomid failed, roomid: %s, err: %s\n", roomId, err)
		return
	}
	room, err := roomBucket.Get(rid(roomIdi))
	if err != nil {
		client.WriteErrorMsg("incorrect roomId.")
		cleaner.CleanClient(client)
		log.Printf("Get room failed, roomid: %s, err: %s\n", roomId, err)
		return
	}
	client.RoomId = rid(roomIdi)
	room.AddClient(client)

	clientBucket.Add(client)

	//send
	go listen(client)

	kali, err := strconv.Atoi(Conf.GetConfig("sys", "keepalive_timeout"))
	if err != nil {
		log.Println(err)
		return
	}
	if kali > 0 {
		keepAlive(client, time.Duration(kali)*time.Second)
	}
}

func listen(client *Client) {
	defer cleaner.CleanClient(client)
	for {
		msgs := Proto{}
		err := client.Conn.ReadJSON(&msgs)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				return
			} else {
				log.Println(err)
				return
			}
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
			err := c.WriteControl(websocket.PingMessage, []byte("keepalive"), time.Now().Add(writeWait))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				log.Println("Ping pong timeout, close client", c)
				cleaner.CleanClient(c)
				return
			}
		}
	}()
}

func messagePusher() {
	for {
		proto := <-broadcast
		roomId := proto.RoomId
		room, err := roomBucket.Get(rid(roomId))
		if err != nil {
			log.Println(err)
			continue
		}
		for _, client := range room.GetClients() {
			client.Write(proto)
		}
	}
}

//TODO 支持CloseHandler

func StartServer() {
	var (
		err error
	)
	if err = InitConfig(); err != nil {
		log.Fatal(err)
	}

	if err = InitRoomBucket(); err != nil {
		log.Fatal(err)
	}

	if err = InitClientBucket(); err != nil {
		log.Fatal(err)
	}

	if err = InitCleaner(); err != nil {
		log.Fatal(err)
	}

	// http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/ws", onConnect)

	go messagePusher()

	addr := ":" + Conf.GetConfig("sys", "port")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
