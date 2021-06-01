package check

import (
	"context"
	"github.com/zhaohaiyu1996/akit/aerrors"
	"github.com/zhaohaiyu1996/akit/check"
	"github.com/zhaohaiyu1996/akit/middleware"
)

func Check() middleware.MiddleWare {
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(c context.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(check.Check); ok {
				if err := v.Check(); err != nil {
					return nil, aerrors.ErrorByMessage(aerrors.ErrorCheck, "check error")
				}
			}
			return next(c, req)
		}
	}
}
