package concurrent

import "sync"

// TODO replace me with generics
type T = interface{}

// Vec minimizes locking times in favor of snapshots and heap allocations. Its usage is thread safe.
// If T is a pointer or points to pointers, the usage is generally
// unsafe. Therefore, ensure that T is a value type in its entirety.
// The following methods are intentionally omitted:
//  - Get, Set: even though these methods may be implemented correctly alone,
//    their usage in e.g. threaded for-loops is generally incorrect and they lead to
//    a false sense of security.
type Vec struct {
	buf   []T
	mutex sync.RWMutex
}

// Copy allocates a new slice of T and returns it.
func (l *Vec) Copy() []T {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.copy()
}

// copy performs the internal copy without locking, because we need it within different locking-contexts.
func (l *Vec) copy() []T {
	cpy := make([]T, len(l.buf), len(l.buf))
	for i := range l.buf {
		cpy[i] = l.buf[i]
	}

	return cpy
}

// Len returns the current length of the underlying list buffer. However, due to concurrency usage this is logically
// just a hint and you cannot expect that e.g. Copy returns this amount of entries.
func (l *Vec) Len() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return len(l.buf)
}

// Drain removes all currently available entries and returns them.
func (l *Vec) Drain() []T {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// allocate a new empty slice, to avoid leaking pointers anyway, so we can just return our internal buffer
	// and forget it.
	tmp := l.buf
	l.buf = make([]T, 0, len(tmp)) //free elements

	return tmp
}

// PushBack appends the given element.
func (l *Vec) PushBack(elem T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buf = append(l.buf, elem)
}

// PopBack removes the last element or returns false.
// Example:
//  var list Vec
//	for val, ok := list.PopBack(); ok; {
//     fmt.Println(val)
//  }
func (l *Vec) PopBack() (T, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if len(l.buf) == 0 {
		return nil, false
	}

	idx := len(l.buf) - 1

	elem := l.buf[idx]
	l.buf[idx] = nil // free element
	l.buf = l.buf[:idx]

	return elem, true
}

// PopFirst removes the first element or returns false. This involves a memcopy of all elements, to ensure
// that we do not get an infinite growing backing array.
// Example:
//  var list Vec
//	for val, ok := list.PopFirst(); ok; {
//     fmt.Println(val)
//  }
func (l *Vec) PopFirst() (T, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if len(l.buf) == 0 {
		return nil, false
	}

	elem := l.buf[0]
	// shift everything left => we could use a slice trick here, but that would cause a memory leak
	copy(l.buf, l.buf[1:])
	l.buf[len(l.buf)-1] = nil // free element
	l.buf = l.buf[:len(l.buf)-1]

	return elem, true
}
