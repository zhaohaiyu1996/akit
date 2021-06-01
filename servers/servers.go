package servers

import "context"

// Server is servers server.
type Server interface {
	Endpoint() (string, error)
	Start() error
	Stop() error
}

// Servers is servers context value.
type Servers struct {
	Kind Kind
}

// Kind defines the type of Transport
type Kind string

// Defines a set of servers kind
const (
	KindARPC  Kind = "gRPC"
	KindAHTTP Kind = "HTTP"
)

type transportKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, sv Servers) context.Context {
	return context.WithValue(ctx, transportKey{}, sv)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (sv Servers, ok bool) {
	sv, ok = ctx.Value(transportKey{}).(Servers)
	return
}
