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

TEXT ·get(SB), NOSPLIT, $0-16
	MOVD	g, R1 // R1 = &g
	MOVD    off+0(FP), R2 // R2 = off
	ADD     R1, R2 // R2 = &g + off = &g[off]
	MOVD	(R2), R1 // R1 = g[off]
	MOVD	R1, ret+8(FP)  // ret R1
	RET

TEXT ·set(SB), NOSPLIT, $0-16
	MOVD	g, R1
	MOVD    off+0(FP), R2
	ADD     R1, R2
	MOVD	ret+8(FP), R1
	MOVD	R1, (R2)
	RET
