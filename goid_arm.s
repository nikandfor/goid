#include "textflag.h"

TEXT ·ID(SB), NOSPLIT, $0-8
	MOVD	0x98(R14), F1 // F1 = g.goid
	MOVD	F1, ret+0(FP) // ret F1
	RET

TEXT ·StartPC(SB), NOSPLIT, $0-4
	MOVW	0x128(R14), R1 // R1 = g.startpc
	MOVW	R1, ret+0(FP)  // ret R1
	RET

TEXT ·GoPC(SB), NOSPLIT, $0-4
	MOVW	0x118(R14), R1 // R1 = g.gopc
	MOVW	R1, ret+0(FP)  // ret R1
	RET
