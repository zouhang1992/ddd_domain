package errors

import "errors"

// 通用领域错误
var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrNotFound       = errors.New("not found")
	ErrCannotDelete   = errors.New("cannot delete")
	ErrInvalidState   = errors.New("invalid state")
)
