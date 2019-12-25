package util

import (
	"sync"
)

//对WaitGroup使用的封装
type WaitGroupWrapper struct {
	sync.WaitGroup
}

//将Add(1)和开始携程，携程结束WaitGroup-1的过程封装成一个函数
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
