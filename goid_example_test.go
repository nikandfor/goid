package goid_test

import (
	"fmt"
	"path"
	"runtime"
	"unsafe"

	"github.com/nikandfor/goid"
)

func ExampleID() {
	c := make(chan struct{})

	loc := func(pc uintptr) string {
		f := runtime.FuncForPC(pc)

		if f == nil {
			return ""
		}

		file, line := f.FileLine(pc)
		file = path.Base(file)

		return fmt.Sprintf("%v:%v", file, line)
	}

	f := func(c chan struct{}) {
		defer close(c)

		id := goid.ID() // this one is simple

		g := goid.GoPC()    // location of "go ...()" instruction
		s := goid.StartPC() // root (the most parent) function entry instruction of that goroutine

		n := runtime.FuncForPC(s).Name()
		n = path.Base(n)

		fmt.Printf("goroutine 0x%x  task %v (%v)  created at %v", id, loc(s), n, loc(g))
	}

	go f(c)

	<-c // wait for f to finish

	// // Output:
	// goroutine 0xc  task goid_example_test.go:19 (goid.ExampleID.func2)  created at goid_example_test.go:35
}

func ExampleGoroutineLocalStorage() {
	type Storage struct {
		Some   int
		Fields string
		Any    map[interface{}]interface{}
	}

	ptr := goid.GLoad()
	fmt.Printf("initial value [%d]: %v\n", goid.ID(), ptr)

	s := &Storage{}

	goid.GSave(unsafe.Pointer(s))

	s.Some = 2
	s.Fields = "can be changed at any time"

	s1 := (*Storage)(goid.GLoad())

	fmt.Printf("loaded [%d]: %+v\n", goid.ID(), s1)

	c := make(chan struct{})
	go func(c chan struct{}) {
		defer close(c)

		s2 := (*Storage)(goid.GLoad())

		fmt.Printf("loaded in another goroutine [%d]: %+v\n", goid.ID(), s2)
	}(c)

	<-c

	// // Output:
	// initial value [1]: <nil>
	// loaded [1]: &{Some:2 Fields:can be changed at any time Any:map[]}
	// loaded in another goroutine [10]: <nil>
}
