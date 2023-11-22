package errors

import (
	"fmt"
	"golang.org/x/text/language"
)

type Error struct {
	Message  string
	HttpCode int
	Code     int

	lang          language.Tag
	originalError error
	translates    map[language.Tag]string
	details       map[string]string
}

func (e *Error) Error() string {
	if v, o := e.translates[e.lang]; o && v != "" {
		return v
	} else {
		fmt.Println("no lang")

		fmt.Println(e.translates)
	}
	return e.Message
}

func (e *Error) AddOriginalError(err error) *Error {
	e.originalError = err
	return e
}

func (e *Error) AddLang(lang string) *Error {
	l, err := language.Parse(lang)
	if err != nil {
		fmt.Println("failed to set lang:", err.Error())
		return e
	}
	e.lang = l
	return e
}

func (e *Error) Original() error {
	if e.originalError != nil {
		return e.originalError
	}
	return nil
}

func (e *Error) AddDetail(key string, detail interface{}) *Error {
	if e.details == nil {
		e.details = make(map[string]string)
	}

	if v, ok := e.details[key]; !ok {
		e.details[key] = fmt.Sprintf("%v", detail)
	} else {
		e.details[key] = fmt.Sprintf("%v - %v", v, detail)
	}

	return e
}

func (e *Error) Details() map[string]string {
	if e.details == nil {
		e.details = make(map[string]string)
	}
	return e.details
}
