package middleware

import (
	"context"
	"fmt"
	"testing"
)

func TestChain(t *testing.T) {
	m := Chain(testPrint("first"), testPrint("second"), testPrint("third"))
	resp, err := m(handle)(context.Background(), "request")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(resp)
}

func handle(ctx context.Context, req interface{}) (interface{}, error) {
	fmt.Println("test middleware chain", req)
	return "handle", nil
}

func testPrint(s string) MiddleWare {
	return func(next MiddleWareFunc) MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			fmt.Println(s + "->start")
			resp, err := next(ctx, req)
			if err != nil {
				return nil, err
			}
			fmt.Println(s + "->end")
			return resp, err
		}
	}
}
