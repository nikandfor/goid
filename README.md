# goid

Do not use this package (https://golang.org/doc/faq#no_goroutine_id).

Or use: example [./goid_example_test.go](./goid_example_test.go), doc https://pkg.go.dev/github.com/nikandfor/goid?tab=doc

tl;dr
```go
func f() {
    id := goid.ID() // <-- that's why you are here

    s := goid.StartPC() // goroutine main (root) function entry point.
    g := goid.GoPC() // go statement pc ("go funcname()")
}
```
