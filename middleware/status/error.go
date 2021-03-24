package status

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/zhaohaiyu1996/akit/aerrors"
	"github.com/zhaohaiyu1996/akit/middleware"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errorConvert is error convert function
type errorConvert func(error) error

// Option is recovery option
type Option func(*options)

// options is errorConvert struct
type options struct {
	handler errorConvert
}

// WithHandler with status handler.
func WithHandler(handler errorConvert) Option {
	return func(o *options) {
		o.handler = handler
	}
}

// NewServerError is a server error middleware
func NewServerError(opts ...Option) middleware.MiddleWare {
	options := options{
		handler: errorEncode,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(err)
			}
			return reply, nil
		}
	}
}

// NewClientError is a client error middleware
func NewClientError(opts ...Option) middleware.MiddleWare {
	options := options{
		handler: errorDecode,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(err)
			}
			return reply, nil
		}
	}
}

// errorDecode is grpc error convert to error interface
func errorDecode(err error) error {
	gs := status.Convert(err)
	se := &aerrors.StatusError{
		Code:    int32(gs.Code()),
		Details: gs.Proto().Details,
	}
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			se.Reason = d.Reason
			se.Message = d.Metadata["message"]
			return se
		}
	}
	return se
}

// errorEncode is error interface convert to grpc error
func errorEncode(err error) error {
	se, ok := aerrors.FromError(err)
	if !ok {
		se = &aerrors.StatusError{
			Code: 2,
		}
	}
	gs := status.Newf(codes.Code(se.Code), "%s: %s", se.Reason, se.Message)
	details := []proto.Message{
		&errdetails.ErrorInfo{
			Reason:   se.Reason,
			Metadata: map[string]string{"message": se.Message},
		},
	}
	for _, any := range se.Details {
		detail := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(any, detail); err != nil {
			continue
		}
		details = append(details, detail.Message)
	}
	gs, err = gs.WithDetails(details...)
	if err != nil {
		return err
	}
	return gs.Err()
}
