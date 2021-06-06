package grpcx

import (
	"context"
	"fmt"
	"github.com/zhaohaiyu1996/akit/internal/host"
	"github.com/zhaohaiyu1996/akit/log"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/servers"
	"google.golang.org/grpc"
	"net"
	"time"
)

const loggerName = "grpcx"

// check *Server is realized by servers.Server
var _ servers.Server = (*Server)(nil)

// Server is a grpcx server wrapper
type Server struct {
	*grpc.Server
	lis        net.Listener
	address    string
	network    string
	log        *log.Logger
	timeout    time.Duration
	middleware middleware.MiddleWare
	grpcOpts   []grpc.ServerOption
}

// NewServer is create a rpc Server
func NewServer(fn func(grpcServer *Server), opts ...ServerOption) *Server {
	var server = &Server{
		address:    ":9426",
		network:    "tcp",
		log:        log.DefaultLogger,
		timeout:    time.Millisecond * 500,
		middleware: middleware.Chain(),
	}
	for _, o := range opts {
		o(server)
	}

	fn(server)

	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptor(server.middleware, server.timeout),
		),
	}

	if len(server.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, server.grpcOpts...)
	}

	server.Server = grpc.NewServer(grpcOpts...)
	server.log.Infof("a rpc server is starting at %s", server.address)
	return server
}

// Start is start Grpc server
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	fmt.Println("start at ", s.address)
	return s.Serve(s.lis)
}

// Stop is Stop Grpc server
func (s *Server) Stop() error {
	s.GracefulStop()
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples: grpcx://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpcx://%s", addr), nil
}

func unaryServerInterceptor(m middleware.MiddleWare, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = servers.NewContext(ctx, servers.Servers{Kind: servers.KindARPC})
		ctx = NewServerContext(ctx, ServerInfo{Server: info.Server, FullMethod: info.FullMethod})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if m != nil {
			h = m(h)
		}
		return h(ctx, req)
	}
}
