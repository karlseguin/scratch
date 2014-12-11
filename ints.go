// pools of arrays for various types
package scratch

import (
	"sync/atomic"
)

var EmptyInts = newStrings(nil, 0)

// Pool of []int
type IntsPool struct {
	misses int64
	size   int
	pool   chan *Ints
}

// Create a pool of []int
func NewInts(size, count int) *IntsPool {
	p := &IntsPool{
		size: size,
		pool: make(chan *Ints, count),
	}

	for i := 0; i < count; i++ {
		p.pool <- newInts(p, size)
	}
	return p
}

// The number of times we tried to checkout from an empty pool
func (p *IntsPool) Misses(reset bool) int64 {
	if reset {
		return atomic.SwapInt64(&p.misses, 0)
	}
	return atomic.LoadInt64(&p.misses)
}

// Get a []int
func (p *IntsPool) Checkout() *Ints {
	select {
	case item := <-p.pool:
		return item
	default:
		atomic.AddInt64(&p.misses, 1)
		return newInts(nil, p.size)
	}
}

type Ints struct {
	length int
	pool   *IntsPool
	values []int
}

func newInts(pool *IntsPool, size int) *Ints {
	return &Ints{
		pool:   pool,
		values: make([]int, size),
	}
}

// Returns false if there's no more room
// The last successfully added item returns false
// (as well all additions thereafter)
func (i *Ints) Add(value int) bool {
	if i.length == len(i.values) {
		return false
	}
	i.values[i.length] = value
	i.length++
	return i.length < len(i.values)
}

// Get the values
func (i *Ints) Values() []int {
	return i.values[:i.length]
}

// Get the values
func (i *Ints) Ids() []int {
	return i.Values()
}

// The number of values
func (i *Ints) Len() int {
	return i.length
}

// Resets to 0 values without releasing back to the pool
func (i *Ints) Reset() {
	i.length = 0
}

// Release the item back to the pool
func (i *Ints) Release() {
	if i.pool != nil {
		i.length = 0
		i.pool.pool <- i
	}
}
