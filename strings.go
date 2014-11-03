// pools of arrays for various types
package scratch

import (
	"sync/atomic"
)

// Pool of []String
type StringsPool struct {
	misses int64
	size   int
	pool   chan *Strings
}

// Create a pool of []String
func NewStrings(size, count int) *StringsPool {
	p := &StringsPool{
		size: size,
		pool: make(chan *Strings, count),
	}

	for i := 0; i < count; i++ {
		p.pool <- newStrings(p, size)
	}
	return p
}

// The number of times we tried to checkout from an empty pool
func (p *StringsPool) Misses() int64 {
	return atomic.LoadInt64(&p.misses)
}

// Get a []String
func (p *StringsPool) Checkout() *Strings {
	select {
	case item := <-p.pool:
		return item
	default:
		atomic.AddInt64(&p.misses, 1)
		return newStrings(nil, p.size)
	}
}

type Strings struct {
	length int
	pool   *StringsPool
	values []string
}

func newStrings(pool *StringsPool, size int) *Strings {
	return &Strings{
		pool:   pool,
		values: make([]string, size),
	}
}

// Returns false if there's no more room
// The last successfully added item returns false
// (as well all additions thereafter)
func (i *Strings) Add(value string) bool {
	if i.length == len(i.values) {
		return false
	}
	i.values[i.length] = value
	i.length++
	return i.length < len(i.values)
}

// Get the values
func (i *Strings) Values() []string {
	return i.values[:i.length]
}

// The number of values
func (i *Strings) Len() int {
	return i.length
}

// The number of vlaues
func (i *Strings) Release() {
	if i.pool != nil {
		i.length = 0
		i.pool.pool <- i
	}
}
