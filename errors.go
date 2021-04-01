package lazlow

import "errors"

var (
	ErrIncompatibleValueType = errors.New("incompatible value type")
	ErrOutOfRange            = errors.New("out of range")
)
