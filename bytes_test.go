package scratch

import (
	. "github.com/karlseguin/expect"
	"reflect"
	"testing"
)

type BytesTests struct{}

func Test_Bytes(t *testing.T) {
	Expectify(new(BytesTests), t)
}

func (_ BytesTests) ItemCapacity() {
	p := NewBytes(9, 1)
	bytes := p.Checkout()
	defer bytes.Release()
	Expect(cap(bytes.Scratch())).To.Equal(9)
}

func (_ BytesTests) PoolCapacity() {
	p := NewBytes(4, 3)
	Expect(cap(p.pool)).To.Equal(3)
}

func (_ BytesTests) CreatesItemOnEmptyPool() {
	p := NewBytes(2, 1)
	bytes1 := p.Checkout()
	bytes2 := p.Checkout()
	Expect(cap(bytes2.Scratch())).To.Equal(2)
	Expect(bytes2.pool).To.Equal(nil)
	bytes1.Release()
	bytes2.Release()
	Expect(len(p.pool)).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (_ BytesTests) ReleasesBackToPool() {
	p := NewBytes(20, 1)
	bytes1 := p.Checkout()
	pointer := reflect.ValueOf(bytes1).Pointer()
	bytes1.Release()

	bytes2 := p.Checkout()
	defer bytes2.Release()
	if reflect.ValueOf(bytes2).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}
