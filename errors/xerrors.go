package errors

import (
	"errors"
	"fmt"
)


var _ error = (*Errors)(nil)

// StatusError contains an error response from the server.
type StatusError = Errors

// Is matches each error in the chain with the target value.
func (e *StatusError) Is(target error) bool {
	err, ok := target.(*StatusError)
	if ok {
		return e.Code == err.Code
	}
	return false
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s details = %+v", e.Code, e.Reason, e.Message, e.Details)
}

// Error returns a Status representing c and msg.
func Error(code int32, reason, message string) error {
	return &StatusError{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

// Errorf returns New(c, fmt.Sprintf(format, a...)).
func Errorf(code int32, reason, format string, a ...interface{}) error {
	return Error(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the status code.
func Code(err error) int32 {
	if err == nil {
		return 0 // ok
	}
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code
	}
	return 2 // unknown
}

// Reason returns the status for a particular error.
// It supports wrapped status.
func Reason(err error) string {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Reason
	}
	return ""
}

// FromError returns status error.
func FromError(err error) (*StatusError, bool) {
	if se := new(StatusError); errors.As(err, &se) {
		return se, true
	}
	return nil, false
}

// As mapping errors.As
func As(err error, target interface{}) interface{} {
	se, ok := target.(*StatusError)
	if !ok {
		return false
	}
	return errors.As(err, &se)
}

// Is mapping status.IS
func Is(err error, statusErr *StatusError) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == statusErr.Code
	}
	return false
}
