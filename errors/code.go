package errors

var NotDefault = &StatusError{
	Code:    404,
	Message: "www",
	Reason:  "wwww",
}
