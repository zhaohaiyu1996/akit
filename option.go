package akit

import (
	"context"
	"github.com/zhaohaiyu1996/akit/transport"
	"os"
)

// Option is an engineOptions option.
type Option func(o *engineOptions)

// engineOptions is an engine options
type engineOptions struct {
	id      string
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	servers []transport.Server
}

// ID with service id.
func WithID(id string) Option {
	return func(o *engineOptions) { o.id = id }
}

// Name with service name.
func WithName(name string) Option {
	return func(o *engineOptions) { o.name = name }
}

// Version with service version.
func WithVersion(version string) Option {
	return func(o *engineOptions) { o.version = version }
}

// Context with service context.
func WithContext(ctx context.Context) Option {
	return func(o *engineOptions) { o.ctx = ctx }
}

// Signal with exit signals.
func WithSignal(sigs ...os.Signal) Option {
	return func(o *engineOptions) { o.sigs = sigs }
}

// Server with transport servers.
func WithServer(srv ...transport.Server) Option {
	return func(o *engineOptions) { o.servers = srv }
}
