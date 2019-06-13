package hllsketch

import (
	"datasketch"
	"fmt"
	"strconv"
	"testing"
)

func TestNewHllSketch(t *testing.T) {
	sk := NewHllSketch(16)
	expected := 100000
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
	}
	uniques := sk.Uniques()
	diff := (uniques-float64(expected))/float64(expected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestHllSketch_Union(t *testing.T) {
	f := NewHllSketchFactory(15)
	sk := f.NewSketch()
	other := f.NewSketch()
	expected := 10000000
	totalExpected := expected * 3 / 2
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
		other.Add("hello world"+strconv.Itoa(i*2))
	}
	newSk, err := sk.Union(other)
	if err != nil {
		t.Error(err)
		return
	}
	uniques := newSk.Uniques()
	diff := (uniques-float64(totalExpected))/float64(totalExpected)
	fmt.Println("union unique id count: ", uniques, ", actual id count: ", totalExpected, ", diff: ", diff)

	skUniques := sk.Uniques()
	diff = (skUniques-float64(expected))/float64(expected)
	fmt.Println("sk_unique id count: ", skUniques, ", actual id count: ", expected, ", diff: ", diff)

	otherUniques := other.Uniques()
	diff = (otherUniques-float64(expected))/float64(expected)
	fmt.Println("other_unique id count: ", otherUniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestThetaSketch_Intersection(t *testing.T) {
	f := NewHllSketchFactory(15)
	sk := f.NewSketch()
	other := f.NewSketch()
	expected := 10000000
	intersecionExpected := expected / 2
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
		other.Add("hello world"+strconv.Itoa(i*2))
	}
	inter, err := datasketch.Intersect(sk, other)
	if err != nil {
		t.Error(err)
		return
	}
	diff := (inter-float64(intersecionExpected))/float64(intersecionExpected)
	fmt.Println("intersect unique id count: ", inter, ", actual id count: ", intersecionExpected, ", diff: ", diff)

	skUniques := sk.Uniques()
	diff = (skUniques-float64(expected))/float64(expected)
	fmt.Println("sk_unique id count: ", skUniques, ", actual id count: ", expected, ", diff: ", diff)

	otherUniques := other.Uniques()
	diff = (otherUniques-float64(expected))/float64(expected)
	fmt.Println("other_unique id count: ", otherUniques, ", actual id count: ", expected, ", diff: ", diff)
}



func Test_kMax(t *testing.T) {
	sk := NewHllSketch(1).(*HllSketch)
	var a uint64 = 0x1
	res := sk.kMax(a)
	if res != 63 {
		fmt.Println("kMax of 1 != 63, actual ", res)
	}
}

func Test_getBucket(t *testing.T) {
	sk := NewHllSketch(8).(*HllSketch)
	var a uint64 = 0x1
	res := sk.getBucket(a<<56)
	if res != 1 {
		fmt.Println("bucket of 4 != 1, actual ", res)
	}
}