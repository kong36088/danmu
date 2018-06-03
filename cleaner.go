package danmu

type Cleaner struct {
}

var (
	cleaner *Cleaner
)

func InitCleaner() error{
	cleaner = new(Cleaner)

	return OK
}

/**
	清除保存的client信息，关闭连接
 */
func (cleaner *Cleaner) CleanClient(client *Client) {
	client.Close()

	room, err := roomBucket.Get(client.RoomId)
	if err == nil {
		delete(room.clients, client.Conn)
	}

	clientBucket.Remove(client.Conn)

	client = nil //for gc
}

//TODO CleanRoom
func (cleaner *Cleaner) CleanRoom(room *Room){

}
