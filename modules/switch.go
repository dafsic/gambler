package modules

import "sync/atomic"

var s int32

func On() {
	atomic.CompareAndSwapInt32(&s, 0, 1)
}

func Off() {
	atomic.CompareAndSwapInt32(&s, 1, 0)
}

func State() int32 {
	return s
}
