package aerrors

import (
	"log"
	"testing"
)

func TestError(t *testing.T) {
	err1 := &StatusError{Code: 404,Reason: "这是一个404",Message: "未找到"}
	err2 := &StatusError{Code: 500}
	res1 := Is(err1,err2)
	if res1 == true {
		log.Fatal("res1 failed")
	}

	err2.Code = 404
	res2 := Is(err1,err2)
	if res2 == false {
		log.Fatal("res2 failed")
	}
	t.Log(Reason(err1))
	t.Log(Code(err1))
	t.Log(err1)
}