package thetasketch

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
)

const UpperBound uint64 = math.MaxUint64
const DefaultPrecision = 4096


type Sketch interface {
	Add(string)
	Uniques() float64
	Union(Sketch) Sketch
	Sub(Sketch) Sketch
	Intersection(Sketch) Sketch
	Bytes() []byte
}


type ThetaSketch struct {
	precision int
	heap    UintHeap
	hashObj hash.Hash64
}

func NewThetaSketch(precision int) Sketch {
	if precision < DefaultPrecision {
		precision = DefaultPrecision
	}

	return &ThetaSketch{precision:precision, heap:NewHeap(precision), hashObj:murmur3.New64()}
}

func NewThetaSketchFromBytes(input []byte) (Sketch, error) {
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

func (sk *ThetaSketch) Union(other Sketch) Sketch {
	var otherTheta *ThetaSketch
	var ok bool
	if otherTheta, ok = other.(*ThetaSketch); !ok {
		return nil
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
	return retSketch
}

func (sk *ThetaSketch) Sub(other Sketch) Sketch {
	var otherTheta *ThetaSketch
	var ok bool
	if otherTheta, ok = other.(*ThetaSketch); !ok {
		return nil
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
	return retSketch
}

func (sk *ThetaSketch) Intersection(other Sketch) Sketch {
	union := sk.Union(other)
	skSub := sk.Sub(other)
	otherSub := other.Sub(sk)
	return union.Sub(skSub).Sub(otherSub)
}

func (sk *ThetaSketch) Bytes() []byte {
	buf := new(bytes.Buffer)
	val := append([]uint64{uint64(sk.precision)}, sk.heap.data...)
	_ = binary.Write(buf, binary.BigEndian, val)
	return buf.Bytes()
}



