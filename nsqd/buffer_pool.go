package nsqd

import (
	"bytes"
	"sync"
)

var bp sync.Pool

//在main之前就会创建对象池
func init() {
	bp.New = func() interface{} {
		return &bytes.Buffer{}
	}
}

func bufferPoolGet() *bytes.Buffer {
	return bp.Get().(*bytes.Buffer)
}

func bufferPoolPut(b *bytes.Buffer) {
	bp.Put(b)
}
