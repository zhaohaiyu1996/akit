package middleware

import (
	"context"
	"encoding/json"
	"github.com/zhaohaiyu1996/akit/log"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/servers"
	"github.com/zhaohaiyu1996/akit/servers/grpcx"
	"time"
)

type Option func(*options)

type options struct {
	logger *log.Logger
}

// WithLogger with middleware logger.
func WithLogger(logger *log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func NewServerLogMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	options.logger.WithStr("middleware", "logging")
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			resp, err = next(ctx, req)
			if sv, ok := servers.FromContext(ctx); ok {
				switch sv.Kind {
				case servers.KindARPC:
					info, ok := grpcx.FromServerContext(ctx)
					if !ok {
						return
					}
					var m = map[string]interface{}{
						"kind":     "server",
						"protocol": "grpcx",
						"cost_us":  time.Since(startTime).Nanoseconds() / 1000,
						"method":   info.FullMethod,
						"req":      req,
					}
					a, err := json.Marshal(m)
					if err == nil {
						options.logger.Info(string(a))
					}
				}
			}

			return
		}
	}
}

func NewClientLogMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	options.logger.WithStr("middleware", "logging")
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			resp, err = next(ctx, req)

			if sv, ok := servers.FromContext(ctx); ok {
				switch sv.Kind {
				case servers.KindARPC:
					info, ok := grpcx.FromServerContext(ctx)
					if !ok {
						return
					}
					var m = map[string]interface{}{
						"kind":     "client",
						"protocol": "grpcx",
						"cost_us":  time.Since(startTime).Nanoseconds() / 1000,
						"method":   info.FullMethod,
						"req":      req,
					}
					a, err := json.Marshal(m)
					if err == nil {
						options.logger.Info(string(a))
					}
				}
			}

			return

		}
	}
}
