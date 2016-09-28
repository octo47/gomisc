package clz

import (
	"math/big"
	"testing"
)

// assembly from math/big/arith_386.s
// TEXT ·bitlenAsm(SB),7,$0
//         BSRL x+0(FP), AX
//         JZ Z1
//         INCL AX
//         MOVL AX, n+4(FP)
//         RET
//
// Z1:     MOVL $0, n+4(FP)
//         RET

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

func TestBitlen(t *testing.T) {
	t.Log(bitlen(17))
}

func bitlenBig(a uint32, buf *big.Int) int {
	buf.SetInt64(int64(a))
	return buf.BitLen()
}

func bitlenAsm(a uint32) int

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
		n = bitlenAsm(x)
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
