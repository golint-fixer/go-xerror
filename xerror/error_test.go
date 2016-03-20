package xerror_test

import (
	"fmt"
	"github.com/ibrt/go-xerror/xerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	err := xerror.New("m2", "m1")
	assert.Equal(t, "m2: m1", err.Error())
	assert.Equal(t, []string{"m2", "m1"}, err.Messages())
	assert.Equal(t, []interface{}{}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap(t *testing.T) {
	err := xerror.Wrap(fmt.Errorf("e1"))
	assert.Equal(t, "e1", err.Error())
	assert.Equal(t, []string{"e1"}, err.Messages())
	assert.Equal(t, []interface{}{}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestIs(t *testing.T) {
	err := xerror.New("m2", "m1")
	assert.True(t, err.Is("m2"))
	assert.False(t, err.Is("m1"))
}

func TestContains(t *testing.T) {
	err := xerror.New("m2", "m1")
	assert.True(t, err.Contains("m2"))
	assert.True(t, err.Contains("m1"))
	assert.False(t, err.Contains("m3"))
}

func TestCopy(t *testing.T) {
	err := xerror.New("m1", "m2")
	cp := err.Copy()
	assert.Equal(t, cp.Error(), err.Error())
	assert.Equal(t, cp.Messages(), err.Messages())
	assert.Equal(t, cp.Debug(), err.Debug())
	assert.Equal(t, cp.Stack(), err.Stack())
}

func TestWithMessages(t *testing.T) {
	err := xerror.New("m1")
	w := err.WithMessages("m3", "m2")
	assert.Equal(t, "m3: m2: m1", w.Error())
	assert.Equal(t, []string{"m3", "m2", "m1"}, w.Messages())
	assert.Equal(t, err.Debug(), w.Debug())
	assert.Equal(t, err.Stack(), w.Stack())
}

func TestWithDebug(t *testing.T) {
	err := xerror.New("m1")
	w := err.WithDebug("d2", "d1")
	assert.Equal(t, err.Error(), w.Error())
	assert.Equal(t, err.Messages(), w.Messages())
	assert.Equal(t, []interface{}{"d2", "d1"}, w.Debug())
	assert.Equal(t, err.Stack(), w.Stack())
}

func TestImplementsError(t *testing.T) {
	var err error
	err = xerror.New("m1")
	assert.Equal(t, "m1", err.Error())
}
