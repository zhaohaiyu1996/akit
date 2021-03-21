package transport

import "context"

// Server is transport server.
type Server interface {
	Endpoint() (string, error)
	Start() error
	Stop() error
}

// Transport is transport context value.
type Transport struct {
	Kind Kind
}

// Kind defines the type of Transport
type Kind string

// Defines a set of transport kind
const (
	KindARPC Kind = "gRPC"
	KindAHTTP Kind = "HTTP"
)

type transportKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, tr Transport) context.Context {
	return context.WithValue(ctx, transportKey{}, tr)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (tr Transport, ok bool) {
	tr, ok = ctx.Value(transportKey{}).(Transport)
	return
}

