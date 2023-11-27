package handlers

import "github.com/h-varmazyar/kiwi/pkg/errors"

var (
	errUnsupportedMessage = errors.NewWithCode("unsupported_message", 1000)
	errInvalidContent     = errors.NewWithCode("invalid_content", 1001)
	errInvalidProxy       = errors.NewWithCode("invalid_proxy", 1002)
)
