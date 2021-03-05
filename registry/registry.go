package registry

import "context"

// service registry interface
type Registry interface {
	// plugins name example: etcd
	Name() string
	// init
	Init(ctx context.Context, opts ...Option) (err error)
	// service registry
	Register(ctx context.Context, service *Service) (err error)
	// Service anti registration
	Unregister(ctx context.Context, service *Service) (err error)
	// service discovery
	GetService(ctx context.Context, name string) (service *Service, err error)
}

// Service
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

// Node is service's node
type Node struct {
	Id     string `json:"id"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

