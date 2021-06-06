package async

import "sync"

type SharedCall interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
	DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
}

type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type SharedGroup struct {
	lock  sync.Mutex
	calls map[string]*call
}

func NewSharedCall() SharedCall {
	return &SharedGroup{
		calls: make(map[string]*call),
	}
}

func (sg *SharedGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	c, done := sg.createCall(key)
	if done {
		return c.value, c.err
	}

	sg.makeCall(c, key, fn)
	return c.value, c.err
}

func (sg *SharedGroup) DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error) {
	c, done := sg.createCall(key)
	if done {
		return c.value, false, c.err
	}

	sg.makeCall(c, key, fn)
	return c.value, true, c.err
}

func (sg *SharedGroup) createCall(key string) (*call, bool) {
	sg.lock.Lock()
	if c, ok := sg.calls[key]; ok {
		sg.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c := new(call)
	c.wg.Add(1)
	sg.calls[key] = c
	sg.lock.Unlock()
	return c, false
}

func (sg *SharedGroup) makeCall(c *call, key string, fn func() (interface{}, error)) {
	defer func() {
		sg.lock.Lock()
		delete(sg.calls, key)
		sg.lock.Unlock()
		c.wg.Done()
	}()

	c.value, c.err = fn()
}
