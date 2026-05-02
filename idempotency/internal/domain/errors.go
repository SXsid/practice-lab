package domain

import "errors"

var ErrInvalidStatusTransition = errors.New("invalid status transition")

var ErrProviderDown = errors.New("payment provider unavailable")
