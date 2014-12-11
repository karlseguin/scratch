// pools of arrays for various types
package scratch

import (
	"sync/atomic"
)

var EmptyBytes = newBytes(nil, 0)

// Pool of []String
type BytesPool struct {
	misses int64
	size   int
	pool   chan *Bytes
}

// Create a pool of []String
func NewBytes(size, count int) *BytesPool {
	p := &BytesPool{
		size: size,
		pool: make(chan *Bytes, count),
	}

	for i := 0; i < count; i++ {
		p.pool <- newBytes(p, size)
	}
	return p
}

// The number of times we tried to checkout from an empty pool
func (p *BytesPool) Misses(reset bool) int64 {
	if reset {
		return atomic.SwapInt64(&p.misses, 0)
	}
	return atomic.LoadInt64(&p.misses)
}

// Get a []String
func (p *BytesPool) Checkout() *Bytes {
	select {
	case item := <-p.pool:
		return item
	default:
		atomic.AddInt64(&p.misses, 1)
		return newBytes(nil, p.size)
	}
}

type Bytes struct {
	length  int
	pool    *BytesPool
	scratch []byte
}

func newBytes(pool *BytesPool, size int) *Bytes {
	return &Bytes{
		pool:    pool,
		scratch: make([]byte, size),
	}
}

func (b *Bytes) Scratch() []byte {
	return b.scratch
}

// The number of vlaues
func (b *Bytes) Release() {
	if b.pool != nil {
		b.length = 0
		b.pool.pool <- b
	}
}
