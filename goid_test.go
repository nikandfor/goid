package goid

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"
)

// A waitReason explains why a goroutine has been stopped.
// See gopark. Do not re-use waitReasons, add new ones.
type waitReason uint8

// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}

type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it's
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn't need write barriers. It's still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr
	pc   uintptr
	g    uintptr
	ctxt unsafe.Pointer
	ret  uintptr
	lr   uintptr
	bp   uintptr // for framepointer-enabled architectures
}

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic    *uintptr // innermost panic - offset known to liblink
	_defer    *uintptr // innermost defer
	m         *uintptr // current m; offset known to arm liblink
	sched     gobuf
	syscallsp uintptr // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc uintptr // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp  uintptr // expected sp at top of stack, to check in traceback
	// param is a generic pointer parameter field used to pass
	// values in particular contexts where other storage for the
	// parameter would be difficult to find. It is currently used
	// in three ways:
	// 1. When a channel operation wakes up a blocked goroutine, it sets param to
	//    point to the sudog of the completed blocking operation.
	// 2. By gcAssistAlloc1 to signal back to its caller that the goroutine completed
	//    the GC cycle. It is unsafe to do so in any other way, because the goroutine's
	//    stack may have moved in the meantime.
	// 3. By debugCallWrap to pass parameters to a new goroutine because allocating a
	//    closure in the runtime is forbidden.
	param        unsafe.Pointer
	atomicstatus uint32
	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid         int64
	schedlink    uintptr
	waitsince    int64      // approx time when the g become blocked
	waitreason   waitReason // if status==Gwaiting

	preempt       bool // preemption signal, duplicates stackguard0 = stackpreempt
	preemptStop   bool // transition to _Gpreempted on preemption; otherwise, just deschedule
	preemptShrink bool // shrink stack at synchronous safe point

	// asyncSafePoint is set if g is stopped at an asynchronous
	// safe point. This means there are frames on the stack
	// without precise pointer information.
	asyncSafePoint bool

	paniconfault bool // panic (instead of crash) on unexpected fault address
	gcscandone   bool // g has scanned stack; protected by _Gscan bit in status
	throwsplit   bool // must not split stack
	// activeStackChans indicates that there are unlocked channels
	// pointing into this goroutine's stack. If true, stack
	// copying needs to acquire channel locks to protect these
	// areas of the stack.
	activeStackChans bool
	// parkingOnChan indicates that the goroutine is about to
	// park on a chansend or chanrecv. Used to signal an unsafe point
	// for stack shrinking. It's a boolean value, but is updated atomically.
	parkingOnChan uint8

	raceignore     int8    // ignore race detection events
	sysblocktraced bool    // StartTrace has emitted EvGoInSyscall about this goroutine
	tracking       bool    // whether we're tracking this G for sched latency statistics
	trackingSeq    uint8   // used to decide whether to track this G
	runnableStamp  int64   // timestamp of when the G last became runnable, only used when tracking
	runnableTime   int64   // the amount of time spent runnable, cleared when running, only used when tracking
	sysexitticks   int64   // cputicks when syscall has returned (for tracing)
	traceseq       uint64  // trace event sequencer
	tracelastp     uintptr // last P emitted an event for this goroutine
	lockedm        uintptr
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
	gopc           uintptr    // pc of go statement that created this goroutine
	ancestors      *[]uintptr // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
	startpc        uintptr    // pc of goroutine function
	racectx        uintptr
	waiting        *uintptr       // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
	cgoCtxt        []uintptr      // cgo traceback context
	labels         unsafe.Pointer // profiler labels
	timer          *uintptr       // cached timer for time.Sleep
	selectDone     uint32         // are we participating in a select and did someone win the race?

	// Per-G GC state

	// gcAssistBytes is this G's GC assist credit in terms of
	// bytes allocated. If this is positive, then the G has credit
	// to allocate gcAssistBytes bytes without assisting. If this
	// is negative, then the G must correct this by performing
	// scan work. We track this in bytes to make it fast to update
	// and check for debt in the malloc hot path. The assist ratio
	// determines how this corresponds to scan work debt.
	gcAssistBytes int64
}

func TestID(t *testing.T) {
	var wg sync.WaitGroup

	par := ID()

	var c0, c1 int64

	wg.Add(1)
	go func() {
		c0 = ID()

		wg.Done()
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		c1 = ID()

		wg.Done()
	}()

	wg.Wait()

	wg.Add(1)
	go testID(t, &wg)

	wg.Wait()

	t.Logf("root %x  c0 %x  c1 %x", par, c0, c1)

	if par == c0 || par == c1 || c0 == c1 {
		t.Errorf("bad ids")
	}

	var g g

	v := reflect.TypeOf(g)

	for _, fn := range []string{"goid", "startpc", "gopc", "ancestors"} {
		f, ok := v.FieldByName(fn)
		if !ok {
			t.Logf("no %v field in g", fn)
			continue
		}
		t.Logf("%-20v: 0x%4x  (%4v)", fn, f.Offset, f.Offset)
	}
}

func testID(t *testing.T, wg *sync.WaitGroup) {
	t.Logf("startpc %v", loc(StartPC()))
	t.Logf("gopc    %v", loc(GoPC()))

	wg.Done()
}

type Storage struct {
	A int
	B string
}

func TestGetSet(t *testing.T) {
	s := &Storage{
		A: 1,
		B: "qweqwe",
	}

	t.Logf("get %p", GLoad())

	GSave(unsafe.Pointer(s))

	t.Logf("set %p <= %p", GLoad(), s)

	q := testGetSet(t)

	if s.B != q {
		t.Errorf("didn't worked")
	}

	t.Logf("res %p", GLoad())
}

func testGetSet(t *testing.T) string {
	time.Sleep(time.Millisecond)

	s := (*Storage)(GLoad())

	return s.B
}

func loc(pc uintptr) string {
	f := runtime.FuncForPC(pc)

	if f == nil {
		return ""
	}

	file, line := f.FileLine(pc)

	file = filepath.Base(file)

	return fmt.Sprintf("%v:%d", file, line)
}

func BenchmarkID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ID()
	}
}
