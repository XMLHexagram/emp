package empErr

import (
	"errors"
	"fmt"
)

type Identifier string

type Error struct {
	Identifier Identifier
	Payload    error
}

func (e *Error) Error() string {
	return fmt.Sprintf("identifier: %s, payload: %s", e.Identifier, e.Payload)
}

func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	err, ok := target.(*Error)
	if !ok {
		return false
	}
	return err.Identifier == e.Identifier
}

func (e *Error) Unwrap() error {
	return e.Payload
}

func (e *Error) Wrap(data interface{}) *Error {
	switch data.(type) {
	case string:
		e.Payload = errors.New(data.(string))
	case error:
		e.Payload = data.(error)
	case []string:
		panic(data)
		e.Payload = errors.New(data.([]string)[0])
	}
	return e
}

func (id Identifier) New() *Error {
	return ErrorMap[id]
}
