package xerror_test

import (
	"errors"
	"fmt"
	"github.com/ibrt/go-xerror/xerror"
	"github.com/stretchr/testify/assert"
	"regexp"
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
	xerr := xerror.Wrap(err)
	assert.Equal(t, err, xerr)
	assert.Nil(t, xerror.Wrap(nil))
}

func TestIs(t *testing.T) {
	err := xerror.New("m2", "m1")
	assert.True(t, err.Is("m2"))
	assert.False(t, err.Is("m1"))
	assert.True(t, xerror.Is(err, "m2"))
	assert.False(t, xerror.Is(err, "m1"))
	assert.True(t, xerror.Is(errors.New("m2"), "m2"))
	assert.False(t, xerror.Is(errors.New("m2"), "m1"))
}

func TestIsPattern(t *testing.T) {
	p1 := regexp.MustCompile("^m")
	p2 := regexp.MustCompile("^m2")
	p3 := regexp.MustCompile("^m1")

	err := xerror.New("m2", "m1")
	assert.True(t, err.IsPattern(p1))
	assert.True(t, err.IsPattern(p2))
	assert.False(t, err.IsPattern(p3))
	assert.True(t, xerror.IsPattern(err, p1))
	assert.True(t, xerror.IsPattern(err, p2))
	assert.False(t, err.IsPattern(p3))
	assert.True(t, xerror.IsPattern(errors.New("m2"), p1))
	assert.True(t, xerror.IsPattern(errors.New("m2"), p2))
	assert.False(t, xerror.IsPattern(errors.New("m2"), p3))
}

func TestContains(t *testing.T) {
	err := xerror.New("m2", "m1")
	assert.True(t, err.Contains("m2"))
	assert.True(t, err.Contains("m1"))
	assert.False(t, err.Contains("m3"))
	assert.True(t, xerror.Contains(err, "m2"))
	assert.True(t, xerror.Contains(err, "m1"))
	assert.False(t, xerror.Contains(err, "m3"))
	assert.True(t, xerror.Contains(errors.New("m2"), "m2"))
	assert.False(t, xerror.Contains(errors.New("m2"), "m1"))
}

func TestContainsPattern(t *testing.T) {
	p1 := regexp.MustCompile("^m")
	p2 := regexp.MustCompile("^m2")
	p3 := regexp.MustCompile("^m1")
	p4 := regexp.MustCompile("^m3")

	err := xerror.New("m2", "m1")
	assert.True(t, err.ContainsPattern(p1))
	assert.True(t, err.ContainsPattern(p2))
	assert.True(t, err.ContainsPattern(p3))
	assert.False(t, err.ContainsPattern(p4))
	assert.True(t, xerror.ContainsPattern(err, p1))
	assert.True(t, xerror.ContainsPattern(err, p2))
	assert.True(t, xerror.ContainsPattern(err, p3))
	assert.False(t, xerror.ContainsPattern(err, p4))
	assert.True(t, xerror.ContainsPattern(errors.New("m2"), p1))
	assert.False(t, xerror.ContainsPattern(errors.New("m2"), p3))
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
	assert.Nil(t, error(xerror.Wrap(nil)))
}
