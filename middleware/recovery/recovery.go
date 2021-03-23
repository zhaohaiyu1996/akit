package recovery

import (
	"context"
	"github.com/zhaohaiyu1996/akit/errors"
	"github.com/zhaohaiyu1996/akit/log"
	"github.com/zhaohaiyu1996/akit/middleware"
	"runtime"
)

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
	logger  log.Logger
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger with recovery logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// NewRecovery is a server middleware that recovers from any panics.
func NewRecovery(opts ...Option) middleware.MiddleWare {
	options := options{
		logger: log.DefaultLogger,
		handler: func(ctx context.Context, req, err interface{}) error {
			return errors.ErrorByMessage(errors.Uncertain, "panic triggered: %v", err)
		},
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/recovery", log.DefaultLogger)
	return func(handler middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					log.Errorf("%v: %+v\n%s\n", rerr, req, buf)

					err = options.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
