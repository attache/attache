package views

import (
	"bytes"
	"sync"
)

var bufferpool = &sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 1024)) },
}

func getbuf() *bytes.Buffer  { return bufferpool.Get().(*bytes.Buffer) }
func putbuf(b *bytes.Buffer) { b.Reset(); bufferpool.Put(b) }

func init() {
	// populate initial buffers
	for i := 0; i < StartingBufferCount; i++ {
		bufferpool.Put(bufferpool.New())
	}
}
