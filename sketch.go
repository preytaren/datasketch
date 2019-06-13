package datasketch

import (
	"fmt"
)

func NotImplementError(msg string) error {
	return fmt.Errorf("method not implemented %s", msg)
}

func NoCompatibleTypeError(t1, t2 string, op string) error {
	return fmt.Errorf("cannot call %s on %s of %s", op, t1, t2)
}

type Sketch interface {
	Add(string)
	Uniques() float64
	Union(Sketch) (Sketch, error)
	Sub(Sketch) (Sketch, error)
	Intersection(Sketch) (Sketch, error)
	Bytes() []byte
	String() string
}

type SketchFactory interface {
	NewSketch() Sketch
}

func Union(sketch, other Sketch) (float64, error) {
	newSk, err := sketch.Union(other)
	if err != nil {
		return 0, err
	}
	return newSk.Uniques(), nil
}

func Intersect(sketch, other Sketch) (float64, error) {
	newSk, err := sketch.Intersection(other)
	if err == nil {
		return newSk.Uniques(), nil
	}

	unions, err := Union(sketch, other)
	if err != nil {
		return 0, err
	}

	return sketch.Uniques()+other.Uniques()-unions, nil
}

