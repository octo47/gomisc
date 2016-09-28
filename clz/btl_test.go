package clz

import (
	"math/big"
	"testing"
)

// exp[(1<<i) % 37] = i
var exp [37]uint8

func init() {
	for i := 0; i <= 32; i++ {
		x := uint64(1) << uint(i)
		exp[x%37] = uint8(i)
	}
}

func bitlen(a uint32) int {
	a |= a >> 1
	a |= a >> 2
	a |= a >> 4
	a |= a >> 8
	a |= a >> 16
	return int(exp[(a%37)+1])
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

func TestBitlen(t *testing.T) {
	bl1 := bitlen(17)
	var buf big.Int
	bl2 := bitlenBig(17, &buf)
	bl3 := bitlen64(17)
	if bl1 != bl2 || bl2 != bl3 {
		t.FailNow()
	}
	t.Log(bl1, bl2, bl3)
}

func bitlenBig(a uint32, buf *big.Int) int {
	buf.SetInt64(int64(a))
	return buf.BitLen()
}

var n int

func benchmark1(b *testing.B, v uint32) {
	x := v
	for i := 0; i < b.N; i++ {
		n = bitlen(x)
	}
}

func benchmark2(b *testing.B, v uint32) {
	x := v
	var buf big.Int
	for i := 0; i < b.N; i++ {
		n = bitlenBig(x, &buf)
	}
}

func benchmark3(b *testing.B, v uint32) {
	x := v
	for i := 0; i < b.N; i++ {
		n = bitlen64(uint64(x))
	}
}

func Benchmark1a(b *testing.B) { benchmark1(b, 3) }
func Benchmark1b(b *testing.B) { benchmark1(b, 17) }
func Benchmark1c(b *testing.B) { benchmark1(b, 3872) }
func Benchmark1d(b *testing.B) { benchmark1(b, 3921486) }

func Benchmark2a(b *testing.B) { benchmark2(b, 3) }
func Benchmark2b(b *testing.B) { benchmark2(b, 17) }
func Benchmark2c(b *testing.B) { benchmark2(b, 3872) }
func Benchmark2d(b *testing.B) { benchmark2(b, 3921486) }

func Benchmark3a(b *testing.B) { benchmark3(b, 3) }
func Benchmark3b(b *testing.B) { benchmark3(b, 17) }
func Benchmark3c(b *testing.B) { benchmark3(b, 3872) }
func Benchmark3d(b *testing.B) { benchmark3(b, 3921486) }
