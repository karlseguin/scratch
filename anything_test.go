package scratch

import (
	. "github.com/karlseguin/expect"
	"reflect"
	"testing"
)

type AnythingTests struct{}

func Test_Anything(t *testing.T) {
	Expectify(new(AnythingTests), t)
}

func (_ AnythingTests) ItemCapacity() {
	p := NewAnything(9, 1)
	ne := p.Checkout()
	defer ne.Release()
	Expect(cap(ne.Values())).To.Equal(9)
}

func (_ AnythingTests) PoolCapacity() {
	p := NewAnything(4, 3)
	Expect(cap(p.pool)).To.Equal(3)
}

func (_ AnythingTests) CreatesItemOnEmptyPool() {
	p := NewAnything(2, 1)
	ne1 := p.Checkout()
	ne2 := p.Checkout()
	Expect(cap(ne2.Values())).To.Equal(2)
	Expect(ne2.pool).To.Equal(nil)
	ne1.Release()
	ne2.Release()
	Expect(len(p.pool)).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (_ AnythingTests) ReleasesBackToPool() {
	p := NewAnything(20, 1)
	ne1 := p.Checkout()
	pointer := reflect.ValueOf(ne1).Pointer()
	ne1.Release()

	ne2 := p.Checkout()
	defer ne2.Release()
	if reflect.ValueOf(ne2).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (_ AnythingTests) AddsValues() {
	ne := newAnything(nil, 3)
	ne.Add(2)
	ne.Add(true)
	ne.Add("spice")
	Expect(ne.Len()).To.Equal(3)
	Expect(ne.Values()).To.Equal([]interface{}{2, true, "spice"})
}

func (_ AnythingTests) SilentlyDropsOverflow() {
	ne := newAnything(nil, 2)
	Expect(ne.Add(2)).To.Equal(true)
	Expect(ne.Add("flow")).To.Equal(false)
	Expect(ne.Add(3.44)).To.Equal(false)
	Expect(ne.Len()).To.Equal(2)
	Expect(ne.Values()).To.Equal([]interface{}{2, "flow"})
}

func (_ AnythingTests) ResetsOnRelease() {
	p := NewAnything(20, 1)
	ne := p.Checkout()
	ne.Add(2)
	ne.Add(3)
	ne.Release()
	Expect(ne.Len()).To.Equal(0)
}

func (_ AnythingTests) Resets() {
	p := NewAnything(20, 1)
	ne := p.Checkout()
	ne.Add(2)
	ne.Add(3)
	ne.Reset()
	Expect(ne.Len()).To.Equal(0)
}
