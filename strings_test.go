package scratch

import (
  "testing"
  "reflect"
  . "github.com/karlseguin/expect"
)

type StringsTests struct{}

func Test_Strings(t *testing.T) {
  Expectify(new(StringsTests), t)
}

func (i *StringsTests) ItemCapacity() {
  p := NewStrings(9, 1)
  strings := p.Checkout()
  defer strings.Release()
  Expect(cap(strings.Values())).To.Equal(9)
}

func (i *StringsTests) PoolCapacity() {
  p := NewStrings(4, 3)
  Expect(cap(p.pool)).To.Equal(3)
}

func (i *StringsTests) CreatesItemOnEmptyPool() {
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

func (i *StringsTests) ReleasesBackToPool() {
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

func (i StringsTests) AddsValues() {
  strings := newStrings(nil, 3)
  strings.Add("b")
  strings.Add("f")
  strings.Add("aa")
  Expect(strings.Len()).To.Equal(3)
  Expect(strings.Values()).To.Equal([]string{"b", "f", "aa"})
}

func (i StringsTests) SilentlyDropsOverflow() {
  strings := newStrings(nil, 2)
  strings.Add("zd")
  strings.Add("4q")
  strings.Add("5541")
  Expect(strings.Len()).To.Equal(2)
  Expect(strings.Values()).To.Equal([]string{"zd", "4q"})
}

func (i StringsTests) ResetsOnRelease() {
  p := NewStrings(20, 1)
  strings := p.Checkout()
  strings.Add("a")
  strings.Add("b")
  strings.Release()
  Expect(strings.Len()).To.Equal(0)
}
