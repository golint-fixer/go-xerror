package xerror

import (
	"fmt"
	"runtime"
)

const (
	maxStackLen = 100
)

func newStack() []string {
	stack := make([]uintptr, maxStackLen)
	stack = stack[:runtime.Callers(2, stack)]
	stackTrace := make([]string, len(stack))

	for i := 0; i < len(stack); i++ {
		file, line := runtime.FuncForPC(stack[i]).FileLine(stack[i] - 1)
		stackTrace[i] = fmt.Sprintf("%s:%d (0x%x)", file, line, stack[i])
	}

	return stackTrace
}
