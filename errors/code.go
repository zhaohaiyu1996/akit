package errors

var NotDefault = &StatusError{
	Code:    404,
	Message: "www",
	Reason:  "wwww",
}

var (
	NotHaveInstance = &StatusError{
		Code: 100001,
		Reason: "not have instance",
	}

)