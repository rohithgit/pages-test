package models

import "fmt"

// CloudError returns a cloud error
type MongoError struct {
	Errnum int
	Errmsg string
}

// Error returns the error string
func (ref MongoError) Error() string {
	return fmt.Sprintf("(%d) %s", ref.Errnum, ref.Errmsg)
}

// NewCloudError returns a new cloud error initialized
// with the error number and error message
func NewError(num int, msg string) *MongoError {
	return &MongoError{
		Errnum: num,
		Errmsg: msg,
	}
}
