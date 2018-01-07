package danmu

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
)

const (
	roomMapCup = 100
)

type (
	rid int64
)
//TODO room auto create
func (r rid) String() string {
	return strconv.Itoa(int(r))
}

var (
	roomList   = []rid{1, 2, 3, 4, 5, 6, 7, 8}
	roomBucket *RoomBucket
)

type Room struct {
	RoomId  rid
	Clients map[*websocket.Conn]*Client
	lck     sync.RWMutex
}

func NewRoom(roomId rid) *Room {
	r := new(Room)
	r.RoomId = roomId
	r.Clients = make(map[*websocket.Conn]*Client)
	return r
}

func (room *Room) AddClient(client *Client) error {
	room.lck.Lock()
	room.Clients[client.Conn] = client
	room.lck.Unlock()

	return nil
}

func (room *Room) GetClients() map[*websocket.Conn]*Client {
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

func InitRoomBucket() error {
	roomBucket = &RoomBucket{
		Rooms: make(map[rid]*Room, roomMapCup),
		lck:   sync.RWMutex{},
	}
	for _, v := range roomList {
		r := NewRoom(v)
		err := roomBucket.Add(r)
		if err != nil {
			log.Exit(err)
		}
	}
	return OK
}

func (rb *RoomBucket) Get(id rid) (*Room, error) {
	rb.lck.RLock()
	defer rb.lck.RUnlock()
	if v, ok := rb.Rooms[id]; ok {
		return v, nil
	} else {
		return nil, ErrRoomDoesNotExist
	}

}

func (rb *RoomBucket) Add(room *Room) error {
	rb.lck.Lock()
	rb.Rooms[room.RoomId] = room
	rb.lck.Unlock()
	return OK
}

func (rb *RoomBucket) Remove(room *Room) error {
	if _, ok := rb.Rooms[room.RoomId]; !ok {
		return ErrRoomDoesNotExist
	}
	rb.lck.Lock()
	delete(rb.Rooms, room.RoomId)
	rb.lck.Unlock()
	return OK
}
