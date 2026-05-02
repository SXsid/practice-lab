package domain

import "errors"

var ErrInvalidStatusTransition = errors.New("invalid status transition")

var ErrNotFound = errors.New("not found")

var ErrCaliming = errors.New("soemelse claimed it alreddy")

var ErrRequestInFlight = errors.New("previous request is still execuitng")

var ErrProviderDown = errors.New("payment provider unavailable")
