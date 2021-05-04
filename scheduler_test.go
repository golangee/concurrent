package concurrent

import (
	"testing"
	"time"
)

func TestNewFixedDelayScheduler(t *testing.T) {
	scheduler := NewFixedDelayScheduler(1*time.Second, func() {
		t.Fatal("should not come here")
	})

	for i:=0;i<10;i++{
		go func() {
			time.Sleep(100)
			scheduler.Stop()
		}()
	}

	time.Sleep(1*time.Second)
}
