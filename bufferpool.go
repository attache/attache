package attache

import (
	"bytes"
	"sync"
)

var bufferpool = &sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, d_BUF_SIZE)) },
}

func getbuf() *bytes.Buffer  { return bufferpool.Get().(*bytes.Buffer) }
func putbuf(b *bytes.Buffer) { b.Reset(); bufferpool.Put(b) }

func init() {
	// populate initial buffers
	for i := 0; i < d_BUF_COUNT; i++ {
		bufferpool.Put(bufferpool.New())
	}
}
