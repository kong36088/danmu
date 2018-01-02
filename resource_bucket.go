package danmu

type ResourceBucket struct {
	Clients map[cid]*Client //cid->Client
	//Clients map[*websocket.Conn]*Client //Conn->Client
	Rooms map[rid]*Room //rid->Room
}
