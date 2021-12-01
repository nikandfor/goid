#include "textflag.h"

TEXT ·ID(SB), NOSPLIT, $0-8
	MOVQ	(TLS), AX     // AX = getg()
	MOVQ	0x98(AX), AX   // AX = AX.goid
	MOVQ	AX, ret+0(FP) // ret = AX
	RET

TEXT ·StartPC(SB), NOSPLIT, $0-8
	MOVQ	(TLS), AX     // AX = getg()
	MOVQ	0x128(AX), AX   // AX = AX.startpc
	MOVQ	AX, ret+0(FP) // ret = AX
	RET

TEXT ·GoPC(SB), NOSPLIT, $0-8
	MOVQ	(TLS), AX     // AX = getg()
	MOVQ	0x118(AX), AX   // AX = AX.gopc
	MOVQ	AX, ret+0(FP) // ret = AX
	RET

TEXT ·get(SB), NOSPLIT, $0-16
	MOVQ	(TLS), AX     // AX = getg()
	MOVQ    off+0(FP), BX
	ADDQ    BX, AX
	MOVQ	(AX), AX   // AX = AX.gopc
	MOVQ	AX, ret+8(FP) // ret = AX
	RET

TEXT ·set(SB), NOSPLIT, $0-16
	MOVQ	(TLS), AX     // AX = getg()
	MOVQ    off+0(FP), BX
	ADDQ    BX, AX
	MOVQ	ret+8(FP), BX // ret = AX
	MOVQ	BX, (AX)   // AX = AX.gopc
	RET
