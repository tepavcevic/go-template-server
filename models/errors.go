package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: email is already in use")
	ErrNotFound   = errors.New("models: no resource found")
)
