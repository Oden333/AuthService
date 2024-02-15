package repository

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user with such email already exists")
	ErrUserNotFound        = errors.New("user doesn't exists")
)
