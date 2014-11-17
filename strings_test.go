package scratch

import (
	. "github.com/karlseguin/expect"
	"reflect"
	"testing"
)

type StringsTests struct{}

func Test_Strings(t *testing.T) {
	Expectify(new(StringsTests), t)
}

func (_ StringsTests) ItemCapacity() {
	p := NewStrings(9, 1)
	strings := p.Checkout()
	defer strings.Release()
	Expect(cap(strings.Values())).To.Equal(9)
}

func (_ StringsTests) PoolCapacity() {
	p := NewStrings(4, 3)
	Expect(cap(p.pool)).To.Equal(3)
}

func (_ StringsTests) CreatesItemOnEmptyPool() {
	p := NewStrings(2, 1)
	strings1 := p.Checkout()
	strings2 := p.Checkout()
	Expect(cap(strings2.Values())).To.Equal(2)
	Expect(strings2.pool).To.Equal(nil)
	strings1.Release()
	strings2.Release()
	Expect(len(p.pool)).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (_ StringsTests) ReleasesBackToPool() {
	p := NewStrings(20, 1)
	strings1 := p.Checkout()
	pointer := reflect.ValueOf(strings1).Pointer()
	strings1.Release()

	strings2 := p.Checkout()
	defer strings2.Release()
	if reflect.ValueOf(strings2).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (_ StringsTests) AddsValues() {
	strings := newStrings(nil, 3)
	strings.Add("b")
	strings.Add("f")
	strings.Add("aa")
	Expect(strings.Len()).To.Equal(3)
	Expect(strings.Values()).To.Equal([]string{"b", "f", "aa"})
}

func (_ StringsTests) SilentlyDropsOverflow() {
	strings := newStrings(nil, 2)
	Expect(strings.Add("zd")).To.Equal(true)
	Expect(strings.Add("4q")).To.Equal(false)
	Expect(strings.Add("5541")).To.Equal(false)
	Expect(strings.Len()).To.Equal(2)
	Expect(strings.Values()).To.Equal([]string{"zd", "4q"})
}

func (_ StringsTests) ResetsOnRelease() {
	p := NewStrings(20, 1)
	strings := p.Checkout()
	strings.Add("a")
	strings.Add("b")
	strings.Release()
	Expect(strings.Len()).To.Equal(0)
}
