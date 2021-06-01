package akit

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/zhaohaiyu1996/akit/registry"
	"github.com/zhaohaiyu1996/akit/servers"
	"github.com/zhaohaiyu1996/akit/servers/arpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type Engine struct {
	id      string
	name    string
	version string
	ctx     context.Context
	cancel  func()

	sigs      []os.Signal
	servers   []servers.Server
	Registrar registry.Registrar

	instance *registry.ServiceInstance
}

// New create an Engine manager.
func New(opts ...Option) *Engine {
	var e = new(Engine)
	e.sigs = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	e.ctx = context.Background()
	if id, err := uuid.NewUUID(); err == nil {
		e.id = id.String()
	}

	for _, o := range opts {
		o(e)
	}

	ctx, cancel := context.WithCancel(e.ctx)
	e.ctx = ctx
	e.cancel = cancel

	return e
}

// Run is run an engine and work
func (e *Engine) Run() error {
	g, ctx := errgroup.WithContext(e.ctx)
	for _, srv := range e.servers {
		srv := srv
		g.Go(func() error {
			// wait for stop signal
			<-ctx.Done()
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}

	if e.Registrar != nil {
		if err := e.Registrar.Register(ctx, e.instance); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, e.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				e.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop is stops the engine.
func (e *Engine) Stop() error {
	if e.cancel != nil {
		e.cancel()
	}
	return nil
}

// NewClientConnect is create a client connect
func NewClientConnect(ctx context.Context, insecure bool, opts ...arpc.ClientOption) (*grpc.ClientConn, error) {
	return arpc.Dial(ctx, insecure, opts...)
}
