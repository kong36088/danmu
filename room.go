package danmu

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	roomMapCup = 100
)

type (
	rid int64
)

func (r rid) String() string {
	return strconv.Itoa(int(r))
}

var (
	roomList   = []rid{1, 2, 3}
	roomBucket *RoomBucket
)

type Room struct {
	RoomId  rid
	Clients map[*Client]*Client
	lck     sync.RWMutex
}

func (room *Room) AddClient(client *Client) error {
	room.lck.Lock()
	room.Clients[client] = client
	room.lck.Unlock()

	return nil
}

func (room *Room) GetClients() map[*Client]*Client{
	room.lck.RLock()
	m := room.Clients
	room.lck.RUnlock()

	return m
}

func (room Room) String() string {
	return fmt.Sprintf("RoomId:%s", room.RoomId)
}

type RoomBucket struct {
	Rooms map[rid]*Room
	lck   sync.RWMutex
}

func InitRoomBucket() {
	roomBucket = &RoomBucket{
		Rooms: make(map[rid]*Room),
		lck:   sync.RWMutex{},
	}
	for _, v := range roomList {
		r := &Room{
			RoomId:  rid(v),
			Clients: make(map[*Client]*Client, roomMapCup),
			lck:     sync.RWMutex{},
		}
		err := roomBucket.AddRoom(r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (rb *RoomBucket) Get(id rid) (*Room, error) {
	if v, ok := rb.Rooms[id]; ok {
		return v, nil
	}else{
		return nil, ErrRoomDoesNotExist
	}

}

func (rb *RoomBucket) AddRoom(room *Room) error {
	rb.lck.Lock()
	rb.Rooms[room.RoomId] = room
	rb.lck.Unlock()
	return OK
}

func (rb *RoomBucket) RemoveRoom(room *Room) error {
	if _, ok := rb.Rooms[room.RoomId]; !ok {
		return ErrRoomDoesNotExist
	}
	delete(rb.Rooms, room.RoomId)
	return OK
}
