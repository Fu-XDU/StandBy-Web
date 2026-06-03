package remote

import "errors"

var (
	ErrUnsupportedPage = errors.New("unsupported page id")
	ErrInvalidConfig   = errors.New("invalid page config")
	ErrInvalidDeviceID = errors.New("invalid device id")
	ErrDeviceNotFound  = errors.New("device not found")
	ErrPageNotFound    = errors.New("page config not found")
	ErrEmptyPatch      = errors.New("no config fields to update")
)
