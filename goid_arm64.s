#include "textflag.h"

TEXT ·ID(SB), NOSPLIT, $0-8
	MOVD	0x98(g), R1 // R1 = g.goid
	MOVD	R1, ret+0(FP) // ret R1
	RET

TEXT ·StartPC(SB), NOSPLIT, $0-8
	MOVD	0x128(g), R1 // R1 = g.startpc
	MOVD	R1, ret+0(FP)  // ret R1
	RET

TEXT ·GoPC(SB), NOSPLIT, $0-8
	MOVD	0x118(g), R1 // R1 = g.gopc
	MOVD	R1, ret+0(FP)  // ret R1
	RET
