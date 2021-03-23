package arpc

import (
	"fmt"
	"github.com/zhaohaiyu1996/akit/internal/host"
	"google.golang.org/grpc"
	"net"
)

type ServerOption func(s *Server)

// Server is a grpc server wrapper
type Server struct {
	*grpc.Server
	lis     net.Listener
	address string
	network string
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func NewServer(opts ...ServerOption) *Server {
	var server = &Server{
		address: ":9426",
		network: "tcp",
	}
	for _, o := range opts {
		o(server)
	}

	server.Server = grpc.NewServer()
	return server
}

// Start is start Grpc server
func (s *Server) Start() error {
	lis,err := net.Listen(s.network,s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	fmt.Println("start at ",s.address)
	return s.Serve(s.lis)
}

// Stop is Stop Grpc server
func (s *Server) Stop() error {
	s.GracefulStop()
	return nil
}

func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}