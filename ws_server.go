package danmu

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	log "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	writeWait   = time.Second
	roomIdFiled = "room"
)

//TODO batch push
//TODO 在线人数
var
(
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}, // 解决域不一致的问题
	} // 将http升级为websocket
)

func onConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	client := NewClient(0, conn)
	roomId := r.FormValue(roomIdFiled)
	roomIdi, err := strconv.Atoi(roomId)
	if err != nil {
		client.WriteErrorMsg("incorrect roomId.")
		cleaner.CleanClient(client)
		log.Error("Parse roomid failed, roomid: %s, err: %s\n", roomId, err)
		return
	}
	room, err := roomBucket.Get(rid(roomIdi))
	if err != nil {
		client.WriteErrorMsg("Room does not exist.")
		cleaner.CleanClient(client)
		log.Error("Room does not exist, roomid: %s, err: %s\n", roomId, err)
		return
	}
	client.RoomId = rid(roomIdi)
	room.AddClient(client)

	clientBucket.Add(client)

	//send
	go listen(client)
	go client.keepAlive()

}

// listen listen message that receive from client
// should be called by goroutines
func listen(client *Client) {
	defer cleaner.CleanClient(client)
	for {
		proto := Proto{}
		err := client.ReadJSON(&proto)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				return
			} else {
				log.Error(err)
				return
			}
		}
		log.Debug(proto)
		switch proto.OP{
		case OPMsg:
			msg := &sarama.ProducerMessage{
				Topic: Topic,
				Value: sarama.ByteEncoder(proto.JsonEncode()),
			}
			producer.Input() <- msg
		case OPCountQuery:
		}

	}
}

func messagePusher() {
	var (
		proto *Proto
	)
	proto = NewProto()
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				//fmt.Printf("%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed

				if err := json.Unmarshal(msg.Value, proto); err != nil {
					log.Error(err)
					continue
				}

				roomId := proto.RoomId
				room, err := roomBucket.Get(rid(roomId))
				if err != nil {
					log.Error(err)
					continue
				}
				log.Debug(proto)
				for _, client := range room.GetClients() {
					client.Write(proto)
				}
			}
		default:
		}
	}

}

func StartServer() {
	var (
		err error
	)
	if err = InitConfig(); err != nil {
		panic(err)
	}

	if err = InitLog(); err != nil {
		panic(err)
	}
	defer CloseLog()
	log.Info("danmu server start")

	kafkaAddrs := Conf.GetConfig("kafka", "address")
	kafkaAddr := strings.Split(kafkaAddrs, ",")

	if err = InitKafka(kafkaAddr); err != nil {
		log.Exit(err)
	}
	defer CloseKafka()

	if err = InitRoomBucket(); err != nil {
		log.Exit(err)
	}

	if err = InitClientBucket(); err != nil {
		log.Exit(err)
	}

	if err = InitCleaner(); err != nil {
		log.Exit(err)
	}

	// http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/ws", onConnect)

	go messagePusher()

	addr := ":" + Conf.GetConfig("sys", "port")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Exit(err)
	}
}
