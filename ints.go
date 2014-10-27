// pools of arrays for various types
package scratch

import (
	"sync/atomic"
)

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
func (p *IntsPool) Misses() int64 {
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
	pool *IntsPool
	values []int
}

func newInts(pool *IntsPool, size int) *Ints {
	return &Ints{
		pool: pool,
		values: make([]int, size),
	}
}

func (i *Ints) Add(value int) {
	if i.length == len(i.values) {
		return
	}
	i.values[i.length] = value
	i.length++
}

func (i *Ints) Values() []int {
	return i.values[:i.length]
}

func (i *Ints) Len() int {
	return i.length
}

func (i *Ints) Release() {
	if i.pool != nil {
		i.length = 0
		i.pool.pool <- i
	}
}
