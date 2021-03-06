package grpcx

import (
	"context"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/registry"
	"github.com/zhaohaiyu1996/akit/servers"
	"github.com/zhaohaiyu1996/akit/servers/grpcx/resolver"
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
	discovery  registry.Discovery
	grpcOpts   []grpc.DialOption
}

// WithClientAddress is set client's address
func WithClientAddress(address string) ClientOption {
	return func(o *Clint) { o.address = address }
}

// WithClientTimeout is set client's timeout
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *Clint) { o.timeout = timeout }
}

// WithClientMiddleware is set client's middleware
func WithClientMiddleware(middlewares ...middleware.MiddleWare) ClientOption {
	return func(o *Clint) { o.middleware = middleware.Chain(middlewares...) }
}

// WithClientGrpcOpts is set client's grpcOpts
func WithClientGrpcOpts(grpcOpts []grpc.DialOption) ClientOption {
	return func(o *Clint) { o.grpcOpts = grpcOpts }
}

// WithDiscovery is set client's discovery
func WithDiscovery(discovery registry.Discovery) ClientOption {
	return func(o *Clint) { o.discovery = discovery }
}

// Dial is dail a connect
func Dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := Clint{
		timeout: 500 * time.Millisecond,
		middleware: middleware.Chain(),
	}
	for _, o := range opts {
		o(&options)
	}

	var grpcOpts = []grpc.DialOption{
		grpc.WithUnaryInterceptor(unaryClientInterceptor(options.middleware, options.timeout)),
	}

	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(resolver.NewBuilder(options.discovery)))
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
		ctx = servers.NewContext(ctx, servers.Servers{Kind: servers.KindARPC})
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
