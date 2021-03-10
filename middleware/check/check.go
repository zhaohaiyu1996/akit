package check

import (
	"context"
	"github.com/zhaohaiyu1996/akit/check"
	"github.com/zhaohaiyu1996/akit/errors"
	"github.com/zhaohaiyu1996/akit/middleware"
)

func Check() middleware.MiddleWare {
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(c context.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(check.Check); ok {
				if err := v.Check(); err != nil {
					return nil, errors.ErrorByMessage(errors.ErrorCheck,"check error")
				}
			}
			return next(c, req)
		}
	}
}
