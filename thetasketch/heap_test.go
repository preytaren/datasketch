package thetasketch

import (
	"fmt"
	"testing"
)

func TestNewHeap(t *testing.T) {
	h := NewHeap(10)
	var i, j uint64
	for i=0; i<20; i++ {
		for j=20; j>0; j-- {
			fmt.Println("insert ", i+j)
			h.Push(i+j)
		}
	}
	fmt.Println(h.data)
}

