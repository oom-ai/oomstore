package database

import "fmt"

const (
	ErrCodeNotFound int = iota + 1
)

var errCodeToMessage = map[int]string{
	ErrCodeNotFound: "not found",
}

type StoreError struct {
	Code       int
	ObjectType string
	ObjectKey  interface{}
}

func (e *StoreError) Error() string {
	return fmt.Sprintf("%s %s, key = %+v", e.ObjectType, errCodeToMessage[e.Code], e.ObjectKey)
}

func NewNotFoundError(objectType string, objectKey interface{}) *StoreError {
	return &StoreError{
		Code:       ErrCodeNotFound,
		ObjectType: objectType,
		ObjectKey:  objectKey,
	}
}

func IsNotFound(err error) bool {
	return isErrCode(err, ErrCodeNotFound)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*StoreError); ok {
		return e.Code == code
	}
	return false
}
