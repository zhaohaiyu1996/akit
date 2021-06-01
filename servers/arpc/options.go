package arpc

import (
	"github.com/zhaohaiyu1996/akit/middleware"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ServerOption func(s *Server)

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

// WithLog is set server's alog
func WithLog(log *log.Logger) ServerOption {
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

func WithRegisterServer(fs ...func(s *Server, v interface{})) ServerOption {
	return func(s *Server) {

	}
}
