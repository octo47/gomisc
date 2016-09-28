#include "textflag.h"

// assembly from math/big/arith_386.s
TEXT ·bitlenAsm(SB), 7, $0
	BSRL x+0(FP), AX
	JZ   Z1
	INCL AX
	MOVL AX, n+4(FP)
	RET

Z1:
	MOVL $0, n+4(FP)
	RET
