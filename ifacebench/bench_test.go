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

func newIterator(size int) iteratorImpl {
	return iteratorImpl{index: 0, size: size}
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
	fixture = make([]Point, 1024*1024)
	for i := range fixture {
		fixture[i] = Point{
			timestamp: uint64(rnd.Int63()),
			value:     rnd.Float64(),
		}
	}
}

var ts uint64
var val float64

func BenchmarkInteface(b *testing.B) {
	for iterator := newIterator(b.N); !iterator.AtEnd(); iterator.Next() {
		ts, val = iterator.Value()
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idx := i % len(fixture)
		ts, val = fixture[idx].timestamp, fixture[idx].value
	}
}
