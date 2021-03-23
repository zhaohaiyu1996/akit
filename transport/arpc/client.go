package arpc

import (
	"context"
	"google.golang.org/grpc"
)

// ClientOption is gRPC client option.
type ClientOption func(o *Clint)

// Clint is gRPC Client
type Clint struct {
	address string
}

func WithDailAddress(address string) ClientOption {
	return func(o *Clint) {
		o.address = address
	}
}

func Dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := Clint{
		address: "127.0.0.1:9426",
	}
	for _, o := range opts {
		o(&options)
	}
	var grpcOpts = []grpc.DialOption{}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	return grpc.DialContext(ctx, options.address, grpcOpts...)
}
