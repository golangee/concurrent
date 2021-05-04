package concurrent

import (
	"fmt"
	"testing"
)

func TestVec_Drain(t *testing.T) {
	var list Vec
	for val, ok := list.PopBack(); ok; {
		fmt.Println(val)
	}
}
