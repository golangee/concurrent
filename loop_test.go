package concurrent

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

func TestForEach(t *testing.T) {
	for i := 0; i < 16; i++ {

		for l := 0; l < 1024; l++ {
			check := make([]bool, l)
			count := int64(0)
			if err := Execute(i, l, nil, func(idx int) error {
				atomic.AddInt64(&count, 1)
				check[idx] = true
				return nil
			}); err != nil {
				t.Fatal(err)
			}

			for x, b := range check {
				if !b {
					t.Fatalf("unvisited index: %v", x)
				}
			}

			if count != int64(l) {
				t.Fatalf("%v: expected %v but got %v\n", i, l, count)
			}
		}
	}

	expectedErr := fmt.Errorf("abc")
	err := Execute(1, 10, nil, func(i int) error {
		if i == 9 {
			return expectedErr
		}

		return nil
	})

	if !errors.Is(err, expectedErr) {
		t.Fatal(err)
	}
}
