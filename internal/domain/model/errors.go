package model

import "errors"

// 领域错误
var (
	ErrNotFound         = errors.New("not found")
	ErrCannotDelete     = errors.New("cannot delete")
	ErrAlreadyExists    = errors.New("already exists")
	ErrRoomNumberExists = errors.New("room number already exists in this location")
	ErrInvalidState     = errors.New("invalid state")
	ErrInvalidCommand   = errors.New("invalid command")
)
