package ratelimit

import (
	"context"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/ratelimit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Option is ratelimit option.
type Option func(*options)

type options struct {
	limit ratelimit.Limiter
}

// set a rate limit method,It must be carried out
func WithLimit(limit ratelimit.Limiter) Option {
	return func(o *options) {
		o.limit = limit
	}
}

// NewRateLimitMiddleware is a RateLimit middleware for server
func NewRateLimitMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.limit == nil {
		panic("WithLimit method not be carried out")
	}
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			allow := options.limit.Allow()
			if !allow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}

			return next(ctx, req)
		}
	}
}