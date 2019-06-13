package thetasketch

import (
	"bytes"
	"datasketch"
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
)

const UpperBound uint64 = math.MaxUint64
const DefaultPrecision = 4096

type ThetaSketch struct {
	precision int
	heap    UintHeap
	hashObj hash.Hash64
}

func NewThetaSketch(precision int) datasketch.Sketch {
	if precision < DefaultPrecision {
		precision = DefaultPrecision
	}

	return &ThetaSketch{precision:precision, heap:NewHeap(precision), hashObj:murmur3.New64()}
}

func NewThetaSketchFromBytes(input []byte) (datasketch.Sketch, error) {
	if len(input) % 8 != 0 {
		return nil, fmt.Errorf("invalid length thetasketch bytes")
	}
	data := make([]uint64, len(input)/8)
	err := binary.Read(bytes.NewBuffer(input), binary.BigEndian, data)
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("invalid format thetasketch bytes %s", err)
	}
	sk := NewThetaSketch(int(data[0])).(*ThetaSketch)
	for _, e := range data[1:] {
		sk.heap.Push(e)
	}
	return sk, nil
}

func (sk *ThetaSketch) hash(key string) uint64 {
	sk.hashObj.Reset()
	_, _ = sk.hashObj.Write([]byte(key))
	return sk.hashObj.Sum64()
}

func (sk *ThetaSketch) Add(key string) {
	hashVal := sk.hash(key)
	sk.heap.Push(hashVal)
}

func (sk *ThetaSketch) Uniques() float64 {
	if !sk.heap.Full() {
		return float64(sk.heap.Len())
	}
	peak := sk.heap.Peak()
	return float64(sk.heap.Len())*(float64(UpperBound)/float64(UpperBound-peak))
}

func (sk *ThetaSketch) Union(other datasketch.Sketch) (datasketch.Sketch, error) {
	var otherTheta *ThetaSketch
	var ok bool
	if otherTheta, ok = other.(*ThetaSketch); !ok {
		return nil, datasketch.NoCompatibleTypeError(sk.String(), other.String(), "union")
	}
	heap := NewHeap(sk.precision)
	for _, n := range sk.heap.Items() {
		heap.Push(n)
	}
	for _, n := range otherTheta.heap.Items() {
		heap.Push(n)
	}
	retSketch := NewThetaSketch(sk.precision).(*ThetaSketch)
	retSketch.heap = heap
	return retSketch, nil
}

func (sk *ThetaSketch) Sub(other datasketch.Sketch) (datasketch.Sketch, error) {
	var otherTheta *ThetaSketch
	var ok bool
	if otherTheta, ok = other.(*ThetaSketch); !ok {
		return nil, datasketch.NoCompatibleTypeError(sk.String(), other.String(), "sub")
	}
	in := make([]uint64, 0)
	for _, n := range sk.heap.Items() {
		if _, ok := otherTheta.heap.dict[n]; ok {
			continue
		}
		in = append(in, n)
	}

	retSketch := NewThetaSketch(len(in)).(*ThetaSketch)
	for _, e := range in {
		retSketch.heap.Push(e)
	}
	return retSketch, nil
}

func (sk *ThetaSketch) Intersection(other datasketch.Sketch) (datasketch.Sketch, error) {
	union, err := sk.Union(other)
	if err != nil {
		return nil, err
	}
	skSub, _ := sk.Sub(other)
	otherSub, _ := other.Sub(sk)
	unionSub, _ := union.Sub(skSub)
	return unionSub.Sub(otherSub)
}

func (sk *ThetaSketch) Bytes() []byte {
	buf := new(bytes.Buffer)
	val := append([]uint64{uint64(sk.precision)}, sk.heap.data...)
	_ = binary.Write(buf, binary.BigEndian, val)
	return buf.Bytes()
}

func (sk *ThetaSketch) String() string {
	return "thetasketch"
}



