package ifacebench

import (
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
	"unsafe"
)

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
	return getTimestamp(idx), getValue(idx)
}

var fixtureSize = 1 << 12
var fixture []byte
var useUnsafe bool

func newFixture(b *testing.B, size int, unsafe bool) {
	fixture = make([]byte, size*16)
	useUnsafe = unsafe
	rnd := rand.New(rand.NewSource(1234))
	flen := fixtureLen()
	tsMaxExpected = 0
	valSumExpected = 0.0
	for i := 0; i < flen; i++ {
		ts := uint64(rnd.Int63())
		if ts > tsMaxExpected {
			tsMaxExpected = ts
		}
		val := rnd.Float64()
		valSumExpected += val
		setTimestamp(i, ts)
		if getTimestamp(i) != ts {
			b.Fatalf("Mismatched stored timestamps %v <> %v", ts, getTimestamp(i))
		}
		setValue(i, val)
		if getValue(i) != val {
			b.Fatalf("Mismatched stored values %v <> %v", ts, getValue(i))
		}
	}
}

func fixtureLen() int {
	return len(fixture) / 16
}

func getTimestamp(idx int) uint64 {
	bidx := idx * 16
	if useUnsafe {
		return *(*uint64)(unsafe.Pointer(&fixture[bidx]))
	}
	return binary.LittleEndian.Uint64(fixture[bidx:])
}

func getValue(idx int) float64 {
	bidx := idx*16 + 8
	if useUnsafe {
		return *(*float64)(unsafe.Pointer(&fixture[bidx]))
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(fixture[bidx:]))
}

func setTimestamp(idx int, val uint64) {
	bidx := idx * 16
	if cap(fixture) < bidx+16 {
		fixture = append(fixture, 16)
	}
	if useUnsafe {
		*(*uint64)(unsafe.Pointer(&fixture[bidx])) = val
	} else {
		binary.LittleEndian.PutUint64(fixture[bidx:], val)
	}
}

func setValue(idx int, val float64) {
	bidx := idx * 16
	if cap(fixture) < bidx+16 {
		fixture = append(fixture, 16)
	}
	bidx += 8
	if useUnsafe {
		*(*uint64)(unsafe.Pointer(&fixture[bidx])) = math.Float64bits(val)
	} else {
		binary.LittleEndian.PutUint64(fixture[bidx:], math.Float64bits(val))
	}
}

var tsMax uint64
var tsMaxExpected uint64
var valSum float64
var valSumExpected float64

func validate(b *testing.B) {
	if tsMax != tsMaxExpected {
		b.Fatalf("Mismatched %v <> %v \n", tsMax, tsMaxExpected)
	}
	if valSum != valSumExpected {
		b.Fatalf("Mismatched %v <> %v \n", valSum, valSumExpected)
	}
}

func BenchmarkIntefaceEncoder(b *testing.B) {
	benchmarkInteface(b, false)
}

func BenchmarkIntefaceUnsafe(b *testing.B) {
	benchmarkInteface(b, true)
}

func benchmarkInteface(b *testing.B, useUnsafe bool) {
	b.StopTimer()
	newFixture(b, fixtureSize, useUnsafe)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for iterator := newIterator(fixtureLen()); !iterator.AtEnd(); iterator.Next() {
			ts, val := iterator.Value()
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
		validate(b)
	}
}

func BenchmarkIntefaceStructEncoder(b *testing.B) {
	benchmarkIntefaceStruct(b, false)
}
func BenchmarkIntefaceStructUnsafe(b *testing.B) {
	benchmarkIntefaceStruct(b, true)
}

func benchmarkIntefaceStruct(b *testing.B, useUnsafe bool) {
	b.StopTimer()
	newFixture(b, fixtureSize, useUnsafe)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		it := newIterator(fixtureLen()).(*iteratorImpl)
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

func BenchmarkDirectEncoder(b *testing.B) {
	benchmarkDirect(b, false)
}

func BenchmarkDirectUnsafe(b *testing.B) {
	benchmarkDirect(b, true)
}

func benchmarkDirect(b *testing.B, useUnsafe bool) {
	b.StopTimer()
	newFixture(b, fixtureSize, useUnsafe)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < fixtureLen(); idx++ {
			ts, val := getTimestamp(idx), getValue(idx)
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
		validate(b)
	}
}

func BenchmarkCallbackEncoder(b *testing.B) {
	benchmarkCallback(b, false)
}
func BenchmarkCallbackUnsafe(b *testing.B) {
	benchmarkCallback(b, true)
}

func benchmarkCallback(b *testing.B, useUnsafe bool) {
	b.StopTimer()
	newFixture(b, fixtureSize, useUnsafe)
	cb := func(ts uint64, val float64) {
		if ts > tsMax {
			tsMax = ts
		}
		valSum += val
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < fixtureLen(); idx++ {
			cb(getTimestamp(idx), getValue(idx))
		}
	}
	validate(b)
}
