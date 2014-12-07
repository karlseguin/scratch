package scratch

import (
	. "github.com/karlseguin/expect"
	"reflect"
	"testing"
)

type IntsTests struct{}

func Test_Ints(t *testing.T) {
	Expectify(new(IntsTests), t)
}

func (_ IntsTests) ItemCapacity() {
	p := NewInts(9, 1)
	ints := p.Checkout()
	defer ints.Release()
	Expect(cap(ints.Values())).To.Equal(9)
}

func (_ IntsTests) PoolCapacity() {
	p := NewInts(4, 3)
	Expect(cap(p.pool)).To.Equal(3)
}

func (_ IntsTests) CreatesItemOnEmptyPool() {
	p := NewInts(2, 1)
	ints1 := p.Checkout()
	ints2 := p.Checkout()
	Expect(cap(ints2.Values())).To.Equal(2)
	Expect(ints2.pool).To.Equal(nil)
	ints1.Release()
	ints2.Release()
	Expect(len(p.pool)).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (_ IntsTests) ReleasesBackToPool() {
	p := NewInts(20, 1)
	ints1 := p.Checkout()
	pointer := reflect.ValueOf(ints1).Pointer()
	ints1.Release()

	ints2 := p.Checkout()
	defer ints2.Release()
	if reflect.ValueOf(ints2).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (_ IntsTests) AddsValues() {
	ints := newInts(nil, 3)
	ints.Add(2)
	ints.Add(9)
	ints.Add(-1)
	Expect(ints.Len()).To.Equal(3)
	Expect(ints.Values()).To.Equal([]int{2, 9, -1})
}

func (_ IntsTests) SilentlyDropsOverflow() {
	ints := newInts(nil, 2)
	Expect(ints.Add(2)).To.Equal(true)
	Expect(ints.Add(9)).To.Equal(false)
	Expect(ints.Add(-1)).To.Equal(false)
	Expect(ints.Len()).To.Equal(2)
	Expect(ints.Values()).To.Equal([]int{2, 9})
}

func (_ IntsTests) ResetsOnRelease() {
	p := NewInts(20, 1)
	ints := p.Checkout()
	ints.Add(2)
	ints.Add(3)
	ints.Release()
	Expect(ints.Len()).To.Equal(0)
}

func (_ IntsTests) Resets() {
	p := NewInts(20, 1)
	ints := p.Checkout()
	ints.Add(2)
	ints.Add(3)
	ints.Reset()
	Expect(ints.Len()).To.Equal(0)
}
