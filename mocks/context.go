package mocks

import time "time"

type Context struct{}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c Context) Done() <-chan struct{} {
	return make(chan struct{})
}

func (c Context) Err() error {
	return nil
}

func (c Context) Value(key interface{}) interface{} {
	return nil
}
