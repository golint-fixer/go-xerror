package xerror

import (
	"fmt"
	"runtime"
)

const (
	maxStackLen = 100
)

func newStack() []string {
	pcs := make([]uintptr, maxStackLen)
	pcs = pcs[:runtime.Callers(2, pcs)]
	stack := make([]string, len(pcs))

	for i := 0; i < len(pcs); i++ {
		file, line := runtime.FuncForPC(pcs[i]).FileLine(pcs[i] - 1)
		stack[i] = fmt.Sprintf("%s:%d (0x%x)", file, line, pcs[i])
	}

	return stack
}
