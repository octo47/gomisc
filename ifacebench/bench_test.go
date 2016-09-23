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
	return getTimestamp(i.index), getValue(i.index)
}

type point struct {
	timestamp uint64
	value     float64
}

const (
	Encoded = 1
	Unsafe  = 2
	Struct  = 3
)

type fixtureDef struct {
	size       int
	encoded    int
	bytesData  []byte
	structData []point
}

var fixtureSize = 1024 * 1024

var fixture = fixtureDef{
	size:       fixtureSize,
	bytesData:  make([]byte, fixtureSize*16),
	structData: make([]point, fixtureSize),
}

func newFixture(b *testing.B, fixtureType int) {
	rnd := rand.New(rand.NewSource(1234))
	fixture.encoded = fixtureType
	tsMaxExpected = 0
	valSumExpected = 0.0
	for i := 0; i < fixture.size; i++ {
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

func getTimestamp(idx int) uint64 {
	switch fixture.encoded {
	case Encoded:
		bidx := idx * 16
		return binary.LittleEndian.Uint64(fixture.bytesData[bidx:])
	case Unsafe:
		bidx := idx * 16
		return *(*uint64)(unsafe.Pointer(&fixture.bytesData[bidx]))
	case Struct:
		return fixture.structData[idx].timestamp
	default:
		panic("Unknown type")
	}
}

func getValue(idx int) float64 {
	switch fixture.encoded {
	case Encoded:
		bidx := idx*16 + 8
		return math.Float64frombits(
			binary.LittleEndian.Uint64(fixture.bytesData[bidx:]))
	case Unsafe:
		bidx := idx*16 + 8
		return *(*float64)(unsafe.Pointer(&fixture.bytesData[bidx]))
	case Struct:
		return fixture.structData[idx].value
	default:
		panic("Unknown type")
	}
}

func setTimestamp(idx int, val uint64) {
	switch fixture.encoded {
	case Encoded:
		bidx := idx * 16
		binary.LittleEndian.PutUint64(fixture.bytesData[bidx:], val)
	case Unsafe:
		bidx := idx * 16
		*(*uint64)(unsafe.Pointer(&fixture.bytesData[bidx])) = val
	case Struct:
		fixture.structData[idx].timestamp = val
	}
}

func setValue(idx int, val float64) {
	switch fixture.encoded {
	case Encoded:
		bidx := idx*16 + 8
		binary.LittleEndian.PutUint64(
			fixture.bytesData[bidx:], math.Float64bits(val))
	case Unsafe:
		bidx := idx*16 + 8
		*(*uint64)(unsafe.Pointer(
			&fixture.bytesData[bidx])) = math.Float64bits(val)
	case Struct:
		fixture.structData[idx].value = val
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
	benchmarkInteface(b, Encoded)
}

func BenchmarkIntefaceUnsafe(b *testing.B) {
	benchmarkInteface(b, Unsafe)
}

func BenchmarkIntefaceStruct(b *testing.B) {
	benchmarkInteface(b, Struct)
}

func benchmarkInteface(b *testing.B, fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for iterator := newIterator(fixture.size); !iterator.AtEnd(); iterator.Next() {
			ts, val := iterator.Value()
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
	}
	validate(b)
}

func BenchmarkIntefaceStructEncoder(b *testing.B) {
	benchmarkIntefaceStruct(b, Encoded)
}
func BenchmarkIntefaceStructUnsafe(b *testing.B) {
	benchmarkIntefaceStruct(b, Unsafe)
}
func BenchmarkIntefaceStructStruct(b *testing.B) {
	benchmarkIntefaceStruct(b, Struct)
}

func benchmarkIntefaceStruct(b *testing.B, fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		it := newIterator(fixture.size).(*iteratorImpl)
		for iterator := it; !iterator.AtEnd(); iterator.Next() {
			ts, val := iterator.Value()
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
	}
	validate(b)
}

func BenchmarkDirectEncoder(b *testing.B) {
	benchmarkDirect(b, Encoded)
}

func BenchmarkDirectUnsafe(b *testing.B) {
	benchmarkDirect(b, Unsafe)
}

func BenchmarkDirectStruct(b *testing.B) {
	benchmarkDirect(b, Struct)
}

func benchmarkDirect(b *testing.B, fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < fixture.size; idx++ {
			ts, val := getTimestamp(idx), getValue(idx)
			if ts > tsMax {
				tsMax = ts
			}
			valSum += val
		}
	}
	validate(b)
}

var cb = func(ts uint64, val float64) {
	if ts > tsMax {
		tsMax = ts
	}
	valSum += val
}

func BenchmarkCallbackEncoder(b *testing.B) {
	benchmarkCallback(b, cb, Encoded)
}
func BenchmarkCallbackUnsafe(b *testing.B) {
	benchmarkCallback(b, cb, Unsafe)
}
func BenchmarkCallbackStruct(b *testing.B) {
	benchmarkCallback(b, cb, Struct)
}

func benchmarkCallback(
	b *testing.B, cb func(ts uint64, val float64), fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		for idx := 0; idx < fixture.size; idx++ {
			cb(getTimestamp(idx), getValue(idx))
		}
	}
	validate(b)
}
