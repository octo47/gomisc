package ifacebench

import (
	"math/rand"
	"testing"
)

type Point struct {
	timestamp uint64
	value     float64
}

type Iterator interface {
	AtEnd() bool
	Next()
	Value() (timestamp uint64, value float64)
}

type iteratorImpl struct {
	index int
	size  int
}

func newIterator(size int) Iterator {
	return &iteratorImpl{index: 0, size: size}
}

func (i *iteratorImpl) AtEnd() bool {
	return i.index >= i.size
}

func (i *iteratorImpl) Next() {
	i.index++
}

func (i *iteratorImpl) Value() (timestamp uint64, value float64) {
	idx := i.index % len(fixture)
	return fixture[idx].timestamp, fixture[idx].value
}

var fixture []Point

func init() {
	rnd := rand.New(rand.NewSource(1234))
	fixture = make([]Point, 1<<6)
	for i := range fixture {
		fixture[i] = Point{
			timestamp: uint64(rnd.Int63()),
			value:     rnd.Float64(),
		}
	}
}

var tsMax uint64
var valSum float64

func validate(b *testing.B) {
	if tsMax != 9159759876629959426 {
		b.FailNow()
	}
	if uint64(valSum) != 30 {
		b.FailNow()
	}
}

func BenchmarkInteface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for iterator := newIterator(len(fixture)); !iterator.AtEnd(); iterator.Next() {
			ts, val := iterator.Value()
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
		validate(b)
	}
}

func BenchmarkIntefaceStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		it := newIterator(len(fixture)).(*iteratorImpl)
		for iterator := it; !iterator.AtEnd(); iterator.Next() {
			ts, val := iterator.Value()
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
		validate(b)
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < len(fixture); idx++ {
			ts, val := fixture[idx].timestamp, fixture[idx].value
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
		validate(b)
	}
}

func BenchmarkCallback(b *testing.B) {
	cb := func(ts uint64, val float64) {
		if ts > tsMax {
			tsMax = ts
		}
		valSum += val
	}
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < len(fixture); idx++ {
			cb(fixture[idx].timestamp, fixture[idx].value)
		}
	}
	validate(b)
}
