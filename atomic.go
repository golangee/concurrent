package concurrent

import "sync/atomic"

// AtomicBool is a lock free type to set and get boolean flags. This needs to be passed as a pointer to make sense.
type AtomicBool struct {
	val int32
}

// Store just overwrites the value.
func (a *AtomicBool) Store(v bool) {
	val := int32(0)
	if v {
		val = 1
	}

	atomic.StoreInt32(&a.val, val)
}

// Load just reads the value. If pointer is nil, always returns false.
func (a *AtomicBool) Load() bool {
	if a == nil {
		return false
	}

	val := atomic.LoadInt32(&a.val)
	if val == 0 {
		return false
	}

	return true
}

// CompareAndSwap executes the compare-and-swap operation.
func (a *AtomicBool) CompareAndSwap(old, new bool) (swapped bool) {
	oldVal := int32(0)
	newVal := int32(0)

	if old {
		oldVal = 1
	}

	if new {
		newVal = 1
	}

	return atomic.CompareAndSwapInt32(&a.val, oldVal, newVal)
}
