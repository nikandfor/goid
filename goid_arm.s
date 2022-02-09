#include "textflag.h"

TEXT ·ID(SB), NOSPLIT, $0-8
	MOVD	0x98(R14), F1 // F1 = g.goid
	MOVD	F1, ret+0(FP) // ret F1
	RET

TEXT ·StartPC(SB), NOSPLIT, $0-4
	MOVW	0x138(R14), R1 // R1 = g.startpc
	MOVW	R1, ret+0(FP)  // ret R1
	RET

TEXT ·GoPC(SB), NOSPLIT, $0-4
	MOVW	0x128(R14), R1 // R1 = g.gopc
	MOVW	R1, ret+0(FP)  // ret R1
	RET

TEXT ·get(SB), NOSPLIT, $0-8
	MOVW	g, R1 // R1 = &g
	MOVW    off+0(FP), R2 // R2 = off
	ADD     R1, R2 // R2 = &g + off = &g[off]
	MOVW	(R2), R1 // R1 = g[off]
	MOVW	R1, ret+4(FP)  // ret R1
	RET

TEXT ·set(SB), NOSPLIT, $0-8
	MOVW	g, R1
	MOVW    off+0(FP), R2
	ADD     R1, R2
	MOVW	ret+4(FP), R1
	MOVW	R1, (R2)
	RET
