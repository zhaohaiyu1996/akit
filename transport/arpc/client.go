package arpc

import (
	"context"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/middleware/recovery"
	"github.com/zhaohaiyu1996/akit/middleware/status"
	"github.com/zhaohaiyu1996/akit/transport"
	"google.golang.org/grpc"
	"time"
)

// ClientOption is gRPC client option.
type ClientOption func(o *Clint)

// Clint is gRPC Client
type Clint struct {
	address    string
	timeout    time.Duration
	middleware middleware.MiddleWare
	grpcOpts   []grpc.DialOption
}

// WithClientAddress is set client's address
func WithClientAddress(address string) ClientOption {
	return func(o *Clint) {
		o.address = address
	}
}

// WithClientTimeout is set client's timeout
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *Clint) {
		o.timeout = timeout
	}
}

// WithClientMiddleware is set client's middleware
func WithClientMiddleware(middleware middleware.MiddleWare) ClientOption {
	return func(o *Clint) {
		o.middleware = middleware
	}
}

// WithClientGrpcOpts is set client's grpcOpts
func WithClientGrpcOpts(grpcOpts []grpc.DialOption) ClientOption {
	return func(o *Clint) {
		o.grpcOpts = grpcOpts
	}
}

// Dial is dail a connect
func Dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := Clint{
		timeout: 500 * time.Millisecond,
		middleware: middleware.Chain(
			recovery.NewRecovery(),
			status.NewClientError(),
		),
	}
	for _, o := range opts {
		o(&options)
	}

	if options.address == "" {
		panic("Please use aprc.WithDailAddress set client dail address")
	}

	var grpcOpts = []grpc.DialOption{
		grpc.WithUnaryInterceptor(unaryClientInterceptor(options.middleware, options.timeout)),
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.DialContext(ctx, options.address, grpcOpts...)
}

func unaryClientInterceptor(m middleware.MiddleWare, timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindARPC})
		ctx = NewClientContext(ctx, ClientInfo{FullMethod: method})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if m != nil {
			h = m(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
