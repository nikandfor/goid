#include "textflag.h"

TEXT ·ID(SB), NOSPLIT, $0-8
	MOVL	(TLS), AX     // AX = getg()
	MOVD	0x98(AX), X1   // X1 = AX.goid
	MOVD	X1, ret+0(FP) // ret = X1
	RET

TEXT ·StartPC(SB), NOSPLIT, $0-4
	MOVL	(TLS), AX     // AX = getg()
	MOVL	0x128(AX), AX   // AX = AX.startpc
	MOVL	AX, ret+0(FP) // ret = AX
	RET

TEXT ·GoPC(SB), NOSPLIT, $0-4
	MOVL	(TLS), AX     // AX = getg()
	MOVL	0x118(AX), AX   // AX = AX.gopc
	MOVL	AX, ret+0(FP) // ret = AX
	RET
