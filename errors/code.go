package errors

var NotDefault = &StatusError{
	Code:    404,
	Message: "www",
	Reason:  "wwww",
}

var (
	Uncertain = &StatusError{
		Code: 10000,
		Reason: "don't know what the mistake is",
	}
	NotHaveInstance = &StatusError{
		Code: 10001,
		Reason: "not have instance",
	}
	ErrorCheck = &StatusError{
		Code: 10002,
		Reason: "check error after middleware and execute",
	}
)