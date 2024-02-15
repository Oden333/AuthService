package service

import "errors"

var (
	ErrorEmptyName  = errors.New("empty Name param")
	ErrUserNotFound = errors.New("user doesn't exists")
)
