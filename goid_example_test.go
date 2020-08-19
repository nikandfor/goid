package goid

import (
	"fmt"
	"path"
	"runtime"
)

func ExampleID() {
	c := make(chan struct{})

	loc := func(pc uintptr) string {
		file, line := runtime.FuncForPC(pc).FileLine(pc)
		file = path.Base(file)

		return fmt.Sprintf("%v:%v", file, line)
	}

	f := func(c chan struct{}) {
		id := ID() // this simple

		g := GoPC()    // location of "go ...()" instruction
		s := StartPC() // root (the most parent) function entry instruction of that goroutine

		n := runtime.FuncForPC(s).Name()
		n = path.Base(n)

		fmt.Printf("goroutine 0x%x  task %v (%v)  created at %v", id, loc(s), n, loc(g))

		close(c)
	}

	go f(c)

	_, _ = <-c // wait for f to finish

	//	// Output:
	//	// goroutine 0xc  task goid_example_test.go:19 (goid.ExampleID.func2)  created at goid_example_test.go:35
}
