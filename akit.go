package akit

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type Engine struct {
	opts     engineOptions
	ctx      context.Context
	cancel   func()
}

// New create an Engine manager.
func New(opts ...Option) *Engine {
	options := engineOptions{
		ctx:    context.Background(),
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &Engine{
		opts:     options,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (e *Engine) Run() error {
	g, ctx := errgroup.WithContext(e.ctx)
	for _, srv := range e.opts.servers {
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, e.opts.sigs...)
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