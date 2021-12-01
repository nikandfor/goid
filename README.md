# goid

Do not use this package (https://golang.org/doc/faq#no_goroutine_id).

Or use: example [./goid_example_test.go](./goid_example_test.go), doc https://pkg.go.dev/github.com/nikandfor/goid?tab=doc

tl;dr
```go
func f() {
    id := goid.ID() // <-- that's why you are here

    s := goid.StartPC() // goroutine main (root) function entry point ("id := goid.ID()").
    g := goid.GoPC() // go statement pc ("go f()")
}
```

# goroutine local storage

There is some rarely used field in goroutine struct. So it could be used to save some stuff.

```go
type Storage struct {
    // ... any fields ...
}

func f() {
    s := (*Storage)(goid.GLoad()) // nil at first
    
    goid.GSave(unsafe.Pointer(&Storage{ /* ... */}))
    
    // ... later event in some other function ...
    
    s1 := (*Storage)(goid.GLoad())
    // use s1.YourField
}
```

# One more time: it's fun, play with it, but try not to use in real code
