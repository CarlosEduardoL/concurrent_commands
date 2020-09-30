package model

import (
	"sync"
)

type atomicInt struct {
	value int64
	mut   sync.Mutex
}

type AtomicInt interface {
	SetValue(int64)
	GetValue(int64)
}

func (a *atomicInt) SetValue(value int64) {
	a.mut.Lock()
	defer a.mut.Unlock()
	a.value = value
}
func (a *atomicInt) GetValue() int64 {
	a.mut.Lock()
	defer a.mut.Unlock()
	return a.value
}
