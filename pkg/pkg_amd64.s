#include "textflag.h"

TEXT ·main(SB), NOSPLIT, $24-0
    MOVQ $0, a-8*2(SP)
    MOVQ $0, b-8*1(SP)

    MOVQ $10, AX
    MOVQ AX, a-8*2(SP)

    MOVQ AX, 0(SP)
    CALL runntime·printint(SB)
    CALL runntime·printnl(SB)

