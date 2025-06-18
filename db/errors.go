package db

import "errors"

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrDuplicatedRecord = errors.New("record already exists")
	ErrInternalServer   = errors.New("unable to get record from database")
)
