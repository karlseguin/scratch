// pools of arrays for various types
package scratch

import (
	"sync/atomic"
)

// Pool of []int
type AnythingPool struct {
	misses int64
	size   int
	pool   chan *Anything
}

// Create a pool of []int
func NewAnything(size, count int) *AnythingPool {
	p := &AnythingPool{
		size: size,
		pool: make(chan *Anything, count),
	}

	for i := 0; i < count; i++ {
		p.pool <- newAnything(p, size)
	}
	return p
}

// The number of times we tried to checkout from an empty pool
func (p *AnythingPool) Misses(reset bool) int64 {
	if reset {
		return atomic.SwapInt64(&p.misses, 0)
	}
	return atomic.LoadInt64(&p.misses)
}

// Get a []int
func (p *AnythingPool) Checkout() *Anything {
	select {
	case item := <-p.pool:
		return item
	default:
		atomic.AddInt64(&p.misses, 1)
		return newAnything(nil, p.size)
	}
}

type Anything struct {
	length int
	pool   *AnythingPool
	values []interface{}
}

func newAnything(pool *AnythingPool, size int) *Anything {
	return &Anything{
		pool:   pool,
		values: make([]interface{}, size),
	}
}

// Returns false if there's no more room
// The last successfully added item returns false
// (as well all additions thereafter)
func (a *Anything) Add(value interface{}) bool {
	if a.length == len(a.values) {
		return false
	}
	a.values[a.length] = value
	a.length++
	return a.length < len(a.values)
}

// Get the values
func (a *Anything) Values() []interface{} {
	return a.values[:a.length]
}

// The length of values
func (a *Anything) Len() int {
	return a.length
}

// Resets to 0 values without releasing back to the pool
func (a *Anything) Reset() {
	a.length = 0
}

// Release the item back to the pool
func (a *Anything) Release() {
	if a.pool != nil {
		a.length = 0
		a.pool.pool <- a
	}
}
