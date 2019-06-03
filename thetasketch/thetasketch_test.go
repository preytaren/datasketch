package thetasketch

import (
	"fmt"
	"strconv"
	"testing"
)

func TestThetaSketch_NotFull(t *testing.T) {
	sk := NewThetaSketch(65532)
	expected := 100000
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
	}
	uniques := sk.Uniques()
	diff := (uniques-float64(expected))/float64(expected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestNewThetaSketch(t *testing.T) {
	sk := NewThetaSketch(65532)
	expected := 10000
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
	}
	uniques := sk.Uniques()
	diff := (uniques-float64(expected))/float64(expected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestThetaSketch_Union(t *testing.T) {
	sk := NewThetaSketch(65532)
	other := NewThetaSketch(65532)
	expected := 10000000
	totalExpected := expected * 3 / 2
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
		other.Add("hello world"+strconv.Itoa(i*2))
	}
	newSk := sk.Union(other)
	uniques := newSk.Uniques()
	diff := (uniques-float64(totalExpected))/float64(totalExpected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", totalExpected, ", diff: ", diff)

	skUniques := sk.Uniques()
	diff = (skUniques-float64(expected))/float64(expected)
	fmt.Println("sk_unique id count: ", skUniques, ", actual id count: ", expected, ", diff: ", diff)

	otherUniques := other.Uniques()
	diff = (otherUniques-float64(expected))/float64(expected)
	fmt.Println("other_unique id count: ", otherUniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestThetaSketch_Sub(t *testing.T) {
	sk := NewThetaSketch(65532)
	other := NewThetaSketch(65532)
	expected := 10000000
	subExpected := expected / 2
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
		other.Add("hello world"+strconv.Itoa(i*2))
	}
	newSk := sk.Sub(other)
	uniques := newSk.Uniques()
	diff := (uniques-float64(subExpected))/float64(subExpected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", subExpected, ", diff: ", diff)

	skUniques := sk.Uniques()
	diff = (skUniques-float64(expected))/float64(expected)
	fmt.Println("sk_unique id count: ", skUniques, ", actual id count: ", expected, ", diff: ", diff)

	otherUniques := other.Uniques()
	diff = (otherUniques-float64(expected))/float64(expected)
	fmt.Println("other_unique id count: ", otherUniques, ", actual id count: ", expected, ", diff: ", diff)
}

func TestThetaSketch_Intersection(t *testing.T) {
	sk := NewThetaSketch(65532)
	other := NewThetaSketch(65532)
	expected := 10000000
	intersecionExpected := expected / 2
	for i:=0; i<expected; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
		other.Add("hello world"+strconv.Itoa(i*2))
	}
	newSk := sk.Intersection(other)
	uniques := newSk.Uniques()
	diff := (uniques-float64(intersecionExpected))/float64(intersecionExpected)
	fmt.Println("unique id count: ", uniques, ", actual id count: ", intersecionExpected, ", diff: ", diff)

	skUniques := sk.Uniques()
	diff = (skUniques-float64(expected))/float64(expected)
	fmt.Println("sk_unique id count: ", skUniques, ", actual id count: ", expected, ", diff: ", diff)

	otherUniques := other.Uniques()
	diff = (otherUniques-float64(expected))/float64(expected)
	fmt.Println("other_unique id count: ", otherUniques, ", actual id count: ", expected, ", diff: ", diff)
}


func BenchmarkNewThetaSketch(b *testing.B) {
	sk := NewThetaSketch(65532)
	for i:=0; i<b.N; i++ {
		sk.Add("hello world"+strconv.Itoa(i))
	}
}
