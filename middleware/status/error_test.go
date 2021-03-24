package status

import (
	"github.com/zhaohaiyu1996/akit/aerrors"
	"log"
	"testing"
)

func TestError(t *testing.T) {
	err := &errors.StatusError{Code: 404,Message: "未找到",Reason: "404"}
	tmp := errorEncode(err)
	resErr := errorDecode(tmp)
	if !errors.Is(resErr,&errors.StatusError{Code: 404}) {
		log.Fatal("错误")
	}
	t.Log(errors.Code(resErr))
}