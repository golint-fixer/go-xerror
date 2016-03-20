package xerror

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

var (
	frameRegexp = regexp.MustCompile("^([^:]+):([0-9]+) \\(([^\\)]+)\\)")
)

func TestNewStack(t *testing.T) {
	stack := newStack() // Note: this is on line 15 (referenced below).
	frame := frameRegexp.FindAllStringSubmatch(stack[0], -1)[0]
	assert.True(t, strings.HasSuffix(frame[1], "/stack_test.go"))
	assert.Equal(t, "15", frame[2])
}
