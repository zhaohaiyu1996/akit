package akit

import (
	"context"
	"github.com/zhaohaiyu1996/akit/registry"
	"github.com/zhaohaiyu1996/akit/servers"
	"os"
)

// Option is an engineOptions option.
type Option func(o *Engine)

// WithID ID with service id.
func WithID(id string) Option {
	return func(o *Engine) { o.id = id }
}

// WithName Name with service name.
func WithName(name string) Option {
	return func(o *Engine) { o.name = name }
}

// WithVersion Version with service version.
func WithVersion(version string) Option {
	return func(o *Engine) { o.version = version }
}

// WithContext Context with service context.
func WithContext(ctx context.Context) Option {
	return func(o *Engine) { o.ctx = ctx }
}

// WithSignal Signal with exit signals.
func WithSignal(sigs ...os.Signal) Option {
	return func(o *Engine) { o.sigs = sigs }
}

// WithServer Server with servers servers.
func WithServer(srv ...servers.Server) Option {
	return func(o *Engine) { o.servers = srv }
}

// WithRegistrar is with registrar
func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *Engine) { o.Registrar = registrar }
}
