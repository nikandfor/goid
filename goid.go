package goid

// ID returns calling goroutine ID
func ID() int64

// StartPC returns calling goroutine start PC (function entry)
func StartPC() uintptr

// GoPC returns PC of go instruction started calling goroutine
func GoPC() uintptr
