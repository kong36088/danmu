package danmu

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"github.com/kong36088/danmu/utils"
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
	roomList       = []rid{1, 2, 3, 4, 5, 6, 7, 8}
	roomBucket     *RoomBucket
)

type Room struct {
	RoomId    rid
	clients   map[*websocket.Conn]*Client
	lck       sync.RWMutex
	protoList *utils.ConcurrentList //Batch push
}

func NewRoom(roomId rid) *Room {
	r := new(Room)
	r.RoomId = roomId
	r.clients = make(map[*websocket.Conn]*Client)
	r.protoList = utils.NewConcurrentList()
	return r
}

func (room *Room) AddClient(client *Client) error {
	room.lck.Lock()
	room.clients[client.Conn] = client
	room.lck.Unlock()

	return nil
}

func (room *Room) GetClients() map[*websocket.Conn]*Client {
	room.lck.RLock()
	m := room.clients
	room.lck.RUnlock()

	return m
}

func (room Room) String() string {
	return fmt.Sprintf("RoomId:%s", room.RoomId)
}

type RoomBucket struct {
	Rooms     map[rid]*Room
	lck       sync.RWMutex
	observers map[interface{}]interface{}
}

func InitRoomBucket() error {
	roomBucket = &RoomBucket{
		Rooms:     make(map[rid]*Room, roomMapCup),
		lck:       sync.RWMutex{},
		observers: make(map[interface{}]interface{}),
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

	rb.notify(RoomActionAdd, room)

	return OK
}

func (rb *RoomBucket) Remove(room *Room) error {
	if _, ok := rb.Rooms[room.RoomId]; !ok {
		return ErrRoomDoesNotExist
	}
	rb.lck.Lock()
	delete(rb.Rooms, room.RoomId)
	rb.lck.Unlock()

	rb.notify(RoomActionDelete, room)

	return OK
}

func (rb *RoomBucket) AttachObserver(observer *RoomObserverInterface) {
	rb.lck.Lock()
	rb.observers[observer] = observer
	rb.lck.Unlock()
}

func (rb *RoomBucket) DetachObserver(observer *RoomObserverInterface) {
	rb.lck.Lock()
	delete(rb.observers, observer)
	rb.lck.Unlock()
}

func (rb *RoomBucket) notify(action int, room *Room) {
	rb.lck.RLock()
	defer rb.lck.RUnlock()

	for _, ob := range rb.observers {
		ob.(RoomObserverInterface).Update(action, room)
	}
}

const (
	RoomActionAdd    = 1
	RoomActionDelete = 2
)

type RoomObserverInterface interface {
	Update(int, *Room)
}
