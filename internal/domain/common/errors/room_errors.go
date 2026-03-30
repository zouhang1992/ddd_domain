package errors

import "errors"

var (
	// ErrRoomNotFound 房间不存在
	ErrRoomNotFound = errors.New("room not found")
	// ErrRoomNotAvailable 房间不可租
	ErrRoomNotAvailable = errors.New("room not available for lease")
)
