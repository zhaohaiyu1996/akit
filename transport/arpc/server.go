package arpc

import (
	"context"
	"fmt"
	"github.com/zhaohaiyu1996/akit/internal/host"
	"github.com/zhaohaiyu1996/akit/log"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/middleware/recovery"
	"github.com/zhaohaiyu1996/akit/middleware/status"
	"github.com/zhaohaiyu1996/akit/transport"
	"google.golang.org/grpc"
	"net"
	"time"
)

const loggerName = "arpc"

// check *Server is realized by transport.Server
var _ transport.Server = (*Server)(nil)

type ServerOption func(s *Server)

// Server is a grpc server wrapper
type Server struct {
	*grpc.Server
	lis        net.Listener
	address    string
	network    string
	log        *log.Helper
	timeout    time.Duration
	middleware middleware.MiddleWare
	grpcOpts   []grpc.ServerOption
}

// WithAddress is set server's address
func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

// WithNetwork is set server's network
func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// WithLog is set server's log
func WithLog(log *log.Helper) ServerOption {
	return func(s *Server) {
		s.log = log
	}
}

// WithTimeout is set server's Timeout
func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// WithMiddleware is set server's middleware
func WithMiddleware(middleware middleware.MiddleWare) ServerOption {
	return func(s *Server) {
		s.middleware = middleware
	}
}

// WithGrpcOpts is set server's grpcOpts
func WithGrpcOpts(grpcOpts []grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = grpcOpts
	}
}

// NewServer is create a rpc Server
func NewServer(opts ...ServerOption) *Server {
	var server = &Server{
		address: ":9426",
		network: "tcp",
		log:     log.NewHelper(loggerName, log.DefaultLogger),
		timeout: time.Millisecond * 500,
		middleware: middleware.Chain(
			recovery.NewRecovery(),
			status.NewServerError(),
		),
	}
	for _, o := range opts {
		o(server)
	}

	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptor(server.middleware, server.timeout),
		),
	}

	if len(server.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, server.grpcOpts...)
	}

	server.Server = grpc.NewServer(grpcOpts...)
	server.log.Info("a rpc server is starting at ", server.address)
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
// examples: arpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("arpc://%s", addr), nil
}

func unaryServerInterceptor(m middleware.MiddleWare, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindARPC})
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
