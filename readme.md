# Scratch
Pools for []int and []string.

You can use them for anything, but I tend to need these most frequently for
scratch work (short-lived work).

## Usage
Both the `StringsPool` and `IntsPool` have the same interface (except the first deals with strings and the later ints).

First, create the pool by specifying the size of each underlying array and the number of items to keep in the pool:

```go
// create a pool of 128 items, each item can hold up to 20 strings
strPool := scratch.NewStrings(20, 128)
// OR
intPool := scratch.NewInts(20, 128)
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
