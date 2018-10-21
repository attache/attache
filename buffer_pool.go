package attache

import (
	"bytes"
	"sync"
)

func init() {
	// populate initial buffers
	for i := 0; i < bufferDefaultCount; i++ {
		bufferpool.Put(bufferpool.New())
	}
}

const (
	bufferDefaultSize  = 1024
	bufferDefaultCount = 100
)

var bufferpool = &sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, bufferDefaultSize)) },
}

func getbuf() *bytes.Buffer  { return bufferpool.Get().(*bytes.Buffer) }
func putbuf(b *bytes.Buffer) { b.Reset(); bufferpool.Put(b) }
