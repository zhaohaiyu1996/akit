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

func NewServer(opts ...ServerOption) *Server {
	var server = &Server{
		address: ":tcp",
		network: ":9426",
	}
	for _, o := range opts {
		o(server)
	}

	return server
}

// Start is start Grpc server
func (s *Server) Start() error {
	lis,err := net.Listen(s.network,s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	return nil
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