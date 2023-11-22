package errors

import (
	"errors"
	"golang.org/x/text/language"
)

func Cast(err error) *Error {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e
	}
	return &Error{
		Message:       err.Error(),
		HttpCode:      DefaultHttpCode,
		Code:          DefaultCode,
		originalError: err,
		translates:    make(map[language.Tag]string),
	}
}
