package danmu

import "errors"

var (
	OK error = nil
	//room
	RoomDoesNotExist = errors.New("room does not exist")
)
