package ratelimit

type Limiter interface {
	Allow() bool
}
