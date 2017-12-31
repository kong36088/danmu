package danmu

import "errors"

var (
	OK error = nil
	//sys
	ErrParamError = errors.New("incorrect param")

	//room
	ErrRoomDoesNotExist = errors.New("room does not exist")
)
