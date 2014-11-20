# Scratch
Pools for []int, []string and []interface{}

You can use them for anything, but I tend to need these most frequently for
scratch work (short-lived work).

## Usage
Aside from the types they deal with, the `StringsPool`, `IntsPool` and `AnytingPool` have the same interface.

First, create the pool by specifying the size of each underlying array and the number of items to keep in the pool:

```go
// create a pool of 128 items, each item can hold up to 20 strings
strPool := scratch.NewStrings(20, 128)
// OR
intPool := scratch.NewInts(20, 128)
// OR
ne := scratch.NewAnything(20, 128)
```

The pools are thread-safe. You can `Checkout` an item, and then `Release` it:

```go
scratch := intPool.Checkout()
defer scratch.Release()

scratch.Add(3)
scratch.Add(9)
// can also use scratch.Len()
for _, value := range scratch.Values() {
  //
}
```

If you `Add` more values than the underlying array can hold, the item is simply dropped/ignored.

## String.Split

The strings object has a `Split(input, sep string) []string` method which behaves like the standard library's `strings.Split` function. Since the underlying pooled `[]string` is used, it performs slightly faster and allocates considerably less memory.

## Bytes

A `BytesPool` can be created via the `BytesPool(size, count)` function. The `Bytes` object which is pooled behaves differently than the `Strings`, `Ints` and `Anything` objects. All it exposes is a `Scratch()` method which exposes the underlying `[]byte`.

```go
bytesPool := scratch.NewBytes(20, 128)
b := bytesPool.Checkout()
defer b.Release()
bytes := b.Scratch()
//bytes is a []byte of len & cap 20
```

The need for `Bytes` is a fairly special case. The [https://github.com/karlseguin/bytepool](bytepool) package is far more idiomatic / useful for most cases.
