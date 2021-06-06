package errors

import srcerrors "errors"

func Is(err, target error) bool {
	return srcerrors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return srcerrors.As(err, target)
}

func Unwrap(err error) error {
	return srcerrors.Unwrap(err)
}
