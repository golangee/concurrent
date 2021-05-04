# concurrent
Some missing concurrent data structures. These types will be refactored accordingly as soon as generics
become available. 

Why not just always a channel? Depends, but a channel is like a queue and has no additional
possibilities to modify its content besides PushBack and PopFront. It is definitely not a one-size-fits-all
tool, and the other like mutex, atomic et al. are there for a reason.

## Vec

Vec minimizes locking times in favor of snapshots and heap allocations. Its usage is thread safe.
If T is a pointer or points to pointers, the usage is generally
unsafe (depends). Therefore, ensure that T is a value type in its entirety.
The following methods are intentionally omitted:
  - Get, Set: even though these methods may be implemented correctly alone,
    their usage in e.g. threaded for-loops is generally incorrect and they lead to
    a false sense of security.
    

Example usage:
```go
var list Vec
for val, ok := list.PopBack(); ok; {
    fmt.Println(val)
}
```

## Execute

Execute blocks and executes with at most n-routines about length indices and calls f with
the index in no defined order.
The first encountered error is returned and pending routines are tried to get cancelled.

The execution is tried to be aborted, if cancel is set. Why no channel here? Because a channel is actually
more complicated:
  * it is slower, because it uses a full memory barrier (mutex), see src/runtime/chan.go. Nothing you want
    for a huge number with very short execution times (like image processing).
  * it is hard to use right and reason about: lifetime and resource leaks? how many cancel elements to consume? Who
    closes whom?
    
If any error occurs, the execution is tried to be cancelled early and an ExecutionError is returned.

Example:
```go
var cancel AtomicBool 
err := Execute(runtime.NumCPU(), 10, &cancel, func(i int) error {
   // ... 
   return nil
})
```

## FixedDelayScheduler
A FixedDelayScheduler executes within a fixed interval. It delays its firing interval by the actual
callback execution duration, so it is guaranteed that at most only one callback runs at a single point of time
and that the time between scheduled jobs matches the fixed delay.