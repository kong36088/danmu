package danmu

import (
	"fmt"
	"sync"
	"strconv"
	"log"
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
	Clients []Client
	lck     sync.RWMutex
}

func (room Room) String() string {
	return fmt.Sprintf("RoomId:%s", room.RoomId)
}

type RoomBucket struct {
	Rooms map[rid]*Room
}

func InitRoomBucket() {
	roomBucket = &RoomBucket{
		Rooms: make(map[rid]*Room),
	}
	for v := range roomList {
		r := &Room{
			RoomId:  rid(v),
			Clients: make([]Client, roomMapCup),
			lck:     sync.RWMutex{},
		}
		err := roomBucket.AddRoom(r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (rb *RoomBucket) AddRoom(room *Room) error {
	rb.Rooms[room.RoomId] = room
	return OK
}

func (rb *RoomBucket) RemoveRoom(room *Room) error {
	if _, ok := rb.Rooms[room.RoomId]; ok {
		return RoomDoesNotExist
	}
	delete(rb.Rooms, room.RoomId)
	return OK
}

func (rb RoomBucket) String() string {
	rt := "rooms:"
	for _, v := range rb.Rooms {
		rt += v.RoomId.String()
	}
	rt += "\n"
	return rt
}
