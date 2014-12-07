// pools of arrays for various types
package scratch

import (
	"strings"
	"sync/atomic"
)

var EmptyStrings = newStrings(nil, 0)

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
func (s *Strings) Add(value string) bool {
	if s.length == len(s.values) {
		return false
	}
	s.values[s.length] = value
	s.length++
	return s.length < len(s.values)
}

// Get the values
func (s *Strings) Values() []string {
	return s.values[:s.length]
}

// The number of values
func (s *Strings) Len() int {
	return s.length
}

func (s *Strings) Split(input, sep string) []string {
	l, position := len(sep), 0
	for {
		index := strings.Index(input[position:], sep)
		if index == -1 {
			s.Add(input[position:])
			return s.Values()
		}
		if s.Add(input[position:position+index]) == false {
			return s.Values()
		}
		position += index + l
	}
}

// Resets to 0 values without releasing back to the pool
func (s *Strings) Reset() {
	s.length = 0
}

// The number of vlaues
func (s *Strings) Release() {
	if s.pool != nil {
		s.length = 0
		s.pool.pool <- s
	}
}
