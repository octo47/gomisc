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
	toGo       int
	nextOffset int
	nextTs     uint64
	nextValue  float64
}

func newIterator(size int) Iterator {
	return &iteratorImpl{toGo: size}
}

func (i *iteratorImpl) AtEnd() bool {
	return i.toGo < 0
}

func (i *iteratorImpl) Next() {
	if i.toGo > 0 {
		i.nextTs, i.nextValue, i.nextOffset = getTuple(i.nextOffset)
	}
	i.toGo--
}

func (i *iteratorImpl) Value() (uint64, float64) {
	return i.nextTs, i.nextValue
}

type point struct {
	timestamp uint64
	value     float64
}

const (
	Encoded = 1
	Unsafe  = 2
	Struct  = 3
	Leb128  = 4
)

type fixtureDef struct {
	size           int
	encoded        int
	tsMaxExpected  uint64
	valSumExpected float64
	bytesData      []byte
	structData     []point
}

var fixtureSize = 1024 * 1024

var fixture fixtureDef

func newFixture(b *testing.B, fixtureType int) {
	rnd := rand.New(rand.NewSource(1234))
	fixture = fixtureDef{
		encoded:        fixtureType,
		size:           0,
		tsMaxExpected:  0,
		valSumExpected: 0.0,
		bytesData:      make([]byte, 0),
		structData:     make([]point, 0),
	}

	for i := 0; i < fixtureSize; i++ {
		ts := uint64(rnd.Int63())
		val := rnd.Float64()
		offset := addTuple(ts, val)
		rts, rval, _ := getTuple(offset)
		if rts != ts {
			b.Fatalf("Mismatched stored timestamps %v <> %v at %d", rts, ts, offset)
		}
		if rval != val {
			b.Fatalf("Mismatched stored values %v <> %v at %d", rval, val, offset)
		}
	}
}

func getTuple(idx int) (ts uint64, val float64, nextOffset int) {
	switch fixture.encoded {
	case Encoded:
		ts = binary.LittleEndian.Uint64(fixture.bytesData[idx : idx+8])
		val = math.Float64frombits(
			binary.LittleEndian.Uint64(fixture.bytesData[idx+8 : idx+16]))
		nextOffset = idx + 16
	case Unsafe:
		ts = *(*uint64)(unsafe.Pointer(&fixture.bytesData[idx]))
		val = *(*float64)(unsafe.Pointer(&fixture.bytesData[idx+8]))
		nextOffset = idx + 16
	case Struct:
		ts, val = fixture.structData[idx].timestamp, fixture.structData[idx].value
		nextOffset = idx + 1
	case Leb128:
		var len int
		ts, len = decodeLeb128(fixture.bytesData[idx:])
		nextOffset = idx + len
		var v uint64
		v, len = decodeLeb128(fixture.bytesData[idx+len:])
		val = math.Float64frombits(v)
		nextOffset += len
	default:
		panic("Unknown type")
	}
	return
}

func addTuple(ts uint64, val float64) int {
	fixture.size++
	if ts > fixture.tsMaxExpected {
		fixture.tsMaxExpected = ts
	}
	fixture.valSumExpected += val
	var buf [8]byte
	switch fixture.encoded {
	case Encoded:
		offset := len(fixture.bytesData)
		binary.LittleEndian.PutUint64(buf[:], ts)
		fixture.bytesData = append(fixture.bytesData, buf[:]...)
		binary.LittleEndian.PutUint64(
			buf[:], math.Float64bits(val))
		fixture.bytesData = append(fixture.bytesData, buf[:]...)
		return offset
	case Unsafe:
		offset := len(fixture.bytesData)
		*(*uint64)(unsafe.Pointer(&buf[0])) = ts
		fixture.bytesData = append(fixture.bytesData, buf[:]...)
		*(*uint64)(unsafe.Pointer(
			&buf[0])) = math.Float64bits(val)
		fixture.bytesData = append(fixture.bytesData, buf[:]...)
		return offset
	case Struct:
		offset := len(fixture.structData)
		fixture.structData = append(fixture.structData, point{ts, val})
		return offset
	case Leb128:
		offset := len(fixture.bytesData)
		fixture.bytesData = encodeLeb128(fixture.bytesData, ts)
		fixture.bytesData = encodeLeb128(fixture.bytesData, math.Float64bits(val))
		return offset
	default:
		panic("Unknown tuple type")
	}
}

const (
	m1  = uint64(0x5555555555555555)
	m2  = uint64(0x3333333333333333)
	m4  = uint64(0x0F0F0F0F0F0F0F0F)
	h01 = uint64(0x0101010101010101)
)

func bitlen64(a uint64) int {
	a |= a >> 1
	a |= a >> 2
	a |= a >> 4
	a |= a >> 8
	a |= a >> 16
	a |= a >> 32
	return popcount(a)
}

// using mult based hamming code (number of '1' in a bit string)
// https://en.wikipedia.org/wiki/Hamming_weight
func popcount(x uint64) int {
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	return int((x * h01) >> 56)    //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
}

func encodeLeb128(buf []byte, v uint64) []byte {
	var b [9]byte
	b[1] = byte(v)
	b[2] = byte(v >> 8)
	b[3] = byte(v >> 16)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 32)
	b[6] = byte(v >> 40)
	b[7] = byte(v >> 48)
	b[8] = byte(v >> 56)
	var i int
	for i = 8; i > 0; i-- {
		if b[i] != 0 {
			break
		}
	}
	b[0] = byte(i)
	return append(buf, b[:i+1]...)
}

func decodeLeb128(buf []byte) (uint64, int) {
	len := int(buf[0])
	var v uint64
	shift := uint32(0)
	for i := 1; i <= len; i++ {
		v += uint64(buf[i]) << shift
		shift += 8
	}
	return v, len + 1
}

var tsMax uint64
var valSum float64

func validate(b *testing.B) {
	if tsMax != fixture.tsMaxExpected {
		b.Fatalf("Mismatched %v <> %v \n", tsMax, fixture.tsMaxExpected)
	}
	if valSum != fixture.valSumExpected {
		b.Fatalf("Mismatched %v <> %v \n", valSum, fixture.valSumExpected)
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

func BenchmarkIntefaceLeb128(b *testing.B) {
	benchmarkInteface(b, Leb128)
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

func BenchmarkIntefaceStructLeb128(b *testing.B) {
	benchmarkIntefaceStruct(b, Leb128)
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

func BenchmarkDirectLeb128(b *testing.B) {
	benchmarkDirect(b, Leb128)
}

func benchmarkDirect(b *testing.B, fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		nextOffset := 0
		for idx := 0; idx < fixture.size; idx++ {
			ts, val, noff := getTuple(nextOffset)
			nextOffset = noff
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
func BenchmarkCallbackLeb128(b *testing.B) {
	benchmarkCallback(b, cb, Leb128)
}

func benchmarkCallback(
	b *testing.B, cb func(ts uint64, val float64), fixtureType int) {
	b.StopTimer()
	newFixture(b, fixtureType)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tsMax = 0
		valSum = 0.0
		nextOffset := 0
		for idx := 0; idx < fixture.size; idx++ {
			ts, val, noff := getTuple(nextOffset)
			nextOffset = noff
			cb(ts, val)
		}
	}
	validate(b)
}
