package concurrent

import (
	"testing"
	"time"
)

func TestNewFixedDelayScheduler(t *testing.T) {
	scheduler := NewFixedDelayScheduler(1*time.Second, func() {
		t.Fatal("should not come here")
	})

	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(100)
			scheduler.Stop()
		}()
	}

	time.Sleep(1 * time.Second)
}

func TestNewFixedDelayScheduler2(t *testing.T) {
	counter := 0
	var scheduler *FixedDelayScheduler

	scheduler = NewFixedDelayScheduler(100*time.Millisecond, func() {
		counter++
		if counter == 5 {
			scheduler.Stop()
		}
	})

	time.Sleep(1 * time.Second)

	if counter != 5 {
		t.Fatalf("expected 5 but counted %d\n", counter)
	}
}
