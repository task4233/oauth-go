package domain

import (
	"errors"
	"fmt"
)

var (
	ErrIsNil = errors.New("the value is nil")
)

type ErrEmpty struct {
	variableName string
}

func NewErrEmpty(variableName string) ErrEmpty {
	return ErrEmpty{variableName: variableName}
}

func (e ErrEmpty) Error() string {
	return fmt.Sprintf("%s is empty", e.variableName)
}
