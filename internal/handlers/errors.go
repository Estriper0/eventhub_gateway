package handlers

import "errors"

var (
	ErrValidateNotPass = errors.New("did not pass validation")
	ErrNoId            = errors.New("no id in params")
)
