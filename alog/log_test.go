package alog

import "testing"

func TestLog(t *testing.T) {
	log := NewLogger(WithLogLevel("info"))
	log.Trace("trace_id")
	log.Info("ZHAOHAIYU")
	log.Errorf("zhaohaiyu:%d", 18)
}

func TestFileLog(t *testing.T) {
	log := NewLogger(WithLogLevel("debug"), WithFileAndCaller("./", true))
	log.Info("123")
}
