package goid

import "unsafe"

const ancestors = 0x130

// ID returns calling goroutine ID
func ID() int64

// StartPC returns calling goroutine start PC (function entry)
func StartPC() uintptr

// GoPC returns PC of go instruction started calling goroutine
func GoPC() uintptr

func GLoad() unsafe.Pointer {
	return get(ancestors)
}

func GSave(p unsafe.Pointer) {
	set(ancestors, p)
}

func get(off uintptr) unsafe.Pointer

func set(off uintptr, val unsafe.Pointer)
