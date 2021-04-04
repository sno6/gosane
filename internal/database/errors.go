package database

import "errors"

var (
	ErrEntityDoesNotExist = errors.New("entity does not exist")
	ErrEntityExists       = errors.New("entity already exists")
)
