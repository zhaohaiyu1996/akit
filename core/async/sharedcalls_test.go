package async

import (
	"fmt"
	"testing"
	"time"
)

func TestNewShareCall(t *testing.T) {
	g := NewSharedCall()
	for i := 0; i < 100; i++ {
		go func() {
			res, err := g.Do("123", func() (interface{}, error) {
				t.Logf("执行：func")
				time.Sleep(time.Second)
				return "res", nil
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		}()

	}
	time.Sleep(time.Second * 2)
}
