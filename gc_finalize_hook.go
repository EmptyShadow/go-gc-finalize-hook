package gcfinalizehook

import (
	"runtime"
	"sync/atomic"
)

type Hook struct {
	f        func()
	ref      *finalizer
	isClosed int32
}

type finalizer struct {
	parent *Hook
}

func NewHook(f func()) *Hook {
	t := &Hook{
		f: f,
	}

	t.ref = &finalizer{
		parent: t,
	}
	runtime.SetFinalizer(t.ref, handleFinalize)
	t.ref = nil

	return t
}

func (t *Hook) Close() {
	atomic.AddInt32(&t.isClosed, 1)
}

func handleFinalize(f *finalizer) {
	if atomic.LoadInt32(&f.parent.isClosed) > 0 {
		return
	}

	f.parent.f()

	runtime.SetFinalizer(f, handleFinalize)
}
