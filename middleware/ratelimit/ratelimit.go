package ratelimit

import (
	"context"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/ratelimit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRateLimitMiddleware(l ratelimit.Limiter) middleware.MiddleWare {
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			allow := l.Allow()
			if !allow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}

			return next(ctx, req)
		}
	}
}