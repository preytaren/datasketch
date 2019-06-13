package hllsketch

import (
	"datasketch"
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
)

const (
	MaxHllSize = 15
	MinHllSize = 1

	HllBase = 1 << 63
	DefaultBucketValue = ^uint64(0)

	HllAlphaInf = 0.721347520444481703680
)

type HllSketch struct {
	bucketN     int
	buckets     []int
	hashObj     hash.Hash64
}

func NewHllSketch(n uint) datasketch.Sketch {
	if n > MaxHllSize {
		n = MaxHllSize
	} else if n < MinHllSize {
		n = MinHllSize
	}
	size := 1
	for i:=MinHllSize; i<=int(n); i++ {
		size <<= 1
	}
	buckets := make([]int, size)
	return &HllSketch{int(n), buckets, murmur3.New64()}
}

func (hll *HllSketch) Add(key string) {
	hll.hashObj.Reset()
	_, _ = hll.hashObj.Write([]byte(key))
	val := hll.hashObj.Sum64()
	hll.insert(val)
}

func (hll *HllSketch) Uniques() float64 {
	m := float64(len(hll.buckets))
	sum := 0.0
	sums := make([]float64, 0)
	for _, b := range hll.buckets {
		sum += 1.0 / math.Pow(2, float64(b))
		sums = append(sums, math.Pow(2, float64(b)))
	}
	return hll.rate()*m*m/sum
}

func (hll *HllSketch) Union(sketch datasketch.Sketch) (datasketch.Sketch, error) {
	if other, ok := sketch.(*HllSketch); !ok {
		return nil, datasketch.NoCompatibleTypeError(hll.String(), sketch.String(), "union")
	} else {
		if other.bucketN != hll.bucketN {
			return nil, datasketch.NoCompatibleTypeError(hll.String()+":"+string(hll.bucketN), sketch.String()+":"+string(other.bucketN), "union")
		}
		newSk := NewHllSketch(uint(hll.bucketN)).(*HllSketch)
		for i, otherMax := range other.buckets {
			newSk.buckets[i] = max(hll.buckets[i], otherMax)
		}
		return newSk, nil
	}
}

func (hll *HllSketch) Sub(sketch datasketch.Sketch) (datasketch.Sketch, error) {
	return nil, datasketch.NotImplementError("Sub")
}

func (hll *HllSketch) Intersection(sketch datasketch.Sketch) (datasketch.Sketch, error) {
	return nil, datasketch.NotImplementError("Intersection, call datasketch.Intersect instead")
}

func (hll *HllSketch) Bytes() []byte {
	return nil
}

func (sk *HllSketch) String() string {
	return "hllsketch"
}

func (hll *HllSketch) insert(ele uint64) {
	bucket := hll.getBucket(ele)
	kMax := hll.kMax(ele)
	if hll.buckets[bucket] < kMax {
		hll.buckets[bucket] = kMax
	}
}

func (hll *HllSketch) getBucket(ele uint64) int {
	for i:=1; i<=64-hll.bucketN; i++ {
		ele >>= 1
	}
	return int(ele)
}

func (hll *HllSketch) kMax(ele uint64) int {
	for i:=1; i<=hll.bucketN; i++ {
		ele <<= 1
	}
	k := 1
	for ele != 0 {
		if (ele & HllBase) != 0 {
			return k
		} else {
			k += 1
		}
		ele <<= 1
	}
	if k == 1 {
		return 64-hll.bucketN+1
	} else {
		return k
	}
}

func (hll *HllSketch) rate() float64 {

	switch hll.bucketN {
	case 4:
		return 0.673
	case 5:
		return 0.697
	case 6:
		return  0.709
	default:
		return HllAlphaInf
	}
}

func max(i, j int) int {
	if i > j {
		return i
	} else {
		return j
	}
}


type DefaultHllSketchFactory struct {
	hllBucketN  int
}

func (dhf *DefaultHllSketchFactory) NewSketch() datasketch.Sketch {
	return NewHllSketch(uint(dhf.hllBucketN))
}

func NewHllSketchFactory(hll int) datasketch.SketchFactory {
	return &DefaultHllSketchFactory{hll}
}