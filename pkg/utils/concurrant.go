package utils

import "sync"

// @ToDo: Experimental and needs to be tested well.
type Array struct {
	sync.RWMutex
	inner []interface{}
}

func NewArray() *Array {
	return &Array{
		inner: make([]interface{}, 0),
	}
}

func (a *Array) Append(v interface{}) {
	a.RLock()
	defer a.RUnlock()
	a.inner = append(a.inner, v)
}
