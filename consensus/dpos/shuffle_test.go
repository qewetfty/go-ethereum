package dpos

import (
	"fmt"
	"testing"
)

func TestShuffle(t *testing.T) {
	var delegateNumber = 4
	for i := 1;i < 100; i++ {
		fmt.Println(Shuffle(int64(i),delegateNumber))
		if i % delegateNumber == 0 {
			fmt.Println("=======================")
		}
	}
}

