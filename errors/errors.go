package errors

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
)

func (e *Error) Error() string {
	return fmt.Sprintf("error:{code:%d,reason:%s,message:%s,metadata:%v}", e.Code, e.Reason, e.Message, e.Metadata)
}

func New(code int, reason, message string) *Error {
	return &Error{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
}

func Newf(code int, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

func (e *Error) ErrorCode() int {
	return int(e.Code)
}

func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

func Code(err error) int {
	if err == nil {
		return 0
	}
	if se := FromError(err); err != nil {
		return int(se.Code)
	}
	return UnknownCode
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				return New(
					int(gs.Code()),
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata)
			}
		}
	}
	return New(UnknownCode, UnknownReason, err.Error())
}

func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}