package danmu

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"sync"
	"time"
)

type MessageRoomObserver struct{}

var (
	commandChans map[*Room]chan string
	lock         *sync.RWMutex
)

func InitMessageHandler() error {
	commandChans = make(map[*Room]chan string)
	lock = &sync.RWMutex{}
	return OK
}

func (mro *MessageRoomObserver) Update(action int, room *Room) {
	lock.Lock()
	defer lock.Unlock()

	if action == RoomActionAdd {
		commandChans[room] = make(chan string)
		go messagePusher(room, commandChans[room])
	} else if action == RoomActionDelete {
		commandChans[room] <- "stop"
		delete(commandChans, room)
	}

}

func messageHandler() {
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
				room.protoList.PushBack(proto)
				/*
				log.Debug(proto)
				for _, client := range room.GetClients() {
					client.Write(proto)
				}*/
			}
		default:
		}
	}
}

func messagePusher(room *Room, commandChan chan string) {
	for {
		select {
		case command := <-commandChan:
			if command == "stop" {
				return
			}
		default:
			if room.protoList.Len() > 0 {
				protos, err := room.protoList.PopAll().([]*Proto)
				if err == false {
					log.Error("[]*Proto type assertion failed")
					continue
				}
				for _, client := range room.GetClients() {
					client.BatchWrite(protos)
				}
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}