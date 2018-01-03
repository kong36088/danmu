package danmu

type Cleaner struct {
}

var (
	cleaner *Cleaner
)

func InitCleaner() {
	cleaner = new(Cleaner)
}

/**
	清除保存的client信息，关闭连接
 */
func (cleaner *Cleaner) CleanClient(client *Client) {
	client.Close()

	room, err := roomBucket.Get(client.RoomId)
	if err == nil {
		delete(room.Clients, client.Conn)
	}

	clientBucket.Remove(client.Conn)

	client = nil //for gc
}
