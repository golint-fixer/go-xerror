package xerror_test

import (
	"encoding/json"
	"errors"
	"github.com/ibrt/go-xerror/xerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_NoPlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.New("fmt")
	assert.Equal(t, "fmt", err.Error())
	assert.Equal(t, []interface{}{}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestNew_PlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.New("fmt %% %v %v", "p2", "p1")
	assert.Equal(t, "fmt % p2 p1", err.Error())
	assert.Equal(t, []interface{}{"p2", "p1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestNew_NoPlaceholdersAndDebug(t *testing.T) {
	err := xerror.New("fmt", "d2", "d1")
	assert.Equal(t, "fmt", err.Error())
	assert.Equal(t, []interface{}{"d2", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestNew_PlaceholdersAndDebug(t *testing.T) {
	err := xerror.New("fmt %% %v %v", "p2", "p1", "d2", "d1")
	assert.Equal(t, "fmt % p2 p1", err.Error())
	assert.Equal(t, []interface{}{"p2", "p1", "d2", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestNew_NoRequiredPlaceholders(t *testing.T) {
	err := xerror.New("fmt %v")
	assert.Equal(t, "fmt %!v(MISSING)", err.Error())
	assert.Equal(t, []interface{}{}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestNew_NotEnoughRequiredPlaceholders(t *testing.T) {
	err := xerror.New("fmt %v %v", "p1")
	assert.Equal(t, "fmt p1 %!v(MISSING)", err.Error())
	assert.Equal(t, []interface{}{"p1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_NilErr(t *testing.T) {
	assert.Panics(t, func() { xerror.Wrap(nil, "fmt") })
}

func TestWrap_NativeErrNoPlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.Wrap(errors.New("ew"), "fmt")
	assert.Equal(t, "fmt: ew", err.Error())
	assert.Equal(t, []interface{}{}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_NativeErrPlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.Wrap(errors.New("ew"), "fmt %% %v %v", "p2", "p1")
	assert.Equal(t, "fmt % p2 p1: ew", err.Error())
	assert.Equal(t, []interface{}{"p2", "p1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_NativeErrNoPlaceholdersAndDebug(t *testing.T) {
	err := xerror.Wrap(errors.New("ew"), "fmt", "d2", "d1")
	assert.Equal(t, "fmt: ew", err.Error())
	assert.Equal(t, []interface{}{"d2", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_NativeErrPlaceholdersAndDebug(t *testing.T) {
	err := xerror.Wrap(errors.New("ew"), "fmt %% %v %v", "p2", "p1", "d2", "d1")
	assert.Equal(t, "fmt % p2 p1: ew", err.Error())
	assert.Equal(t, []interface{}{"p2", "p1", "d2", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_ErrorNoPlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1", "d1"), "fmt2")
	assert.Equal(t, "fmt2: fmt p1", err.Error())
	assert.Equal(t, []interface{}{"p1", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_ErrorPlaceholdersAndNoDebug(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1", "d1"), "fmt2 %% %v %v", "p3", "p2")
	assert.Equal(t, "fmt2 % p3 p2: fmt p1", err.Error())
	assert.Equal(t, []interface{}{"p3", "p2", "p1", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_ErrorNoPlaceholdersAndDebug(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1", "d1"), "fmt2", "d3", "d2")
	assert.Equal(t, "fmt2: fmt p1", err.Error())
	assert.Equal(t, []interface{}{"d3", "d2", "p1", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestWrap_ErrorPlaceholdersAndDebug(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1", "d1"), "fmt2 %% %v %v", "p3", "p2", "d3", "d2")
	assert.Equal(t, "fmt2 % p3 p2: fmt p1", err.Error())
	assert.Equal(t, []interface{}{"p3", "p2", "d3", "d2", "p1", "d1"}, err.Debug())
	assert.True(t, len(err.Stack()) > 0)
}

func TestIs_Method(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1"), "fmt2 %v", "p2")
	assert.Equal(t, "fmt2 p2: fmt p1", err.Error())
	assert.True(t, err.Is("fmt2 %v"))
	assert.False(t, err.Is("fmt2 p2"))
	assert.False(t, err.Is("fmt %v"))
	assert.False(t, err.Is("fmt p1"))
}

func TestIs_TopLevelNilErr(t *testing.T) {
	assert.False(t, xerror.Is(nil, "msg"))
}

func TestIs_TopLevelNativeErr(t *testing.T) {
	err := errors.New("msg")
	assert.True(t, xerror.Is(err, "msg"))
	assert.False(t, xerror.Is(err, "else"))
}

func TestIsTopLevelError(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1"), "fmt2 %v", "p2")
	assert.Equal(t, "fmt2 p2: fmt p1", err.Error())
	assert.True(t, xerror.Is(err, "fmt2 %v"))
	assert.False(t, xerror.Is(err, "fmt2 p2"))
	assert.False(t, xerror.Is(err, "fmt %v"))
	assert.False(t, xerror.Is(err, "fmt p1"))
}

func TestContains_Method(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1"), "fmt2 %v", "p2")
	assert.Equal(t, "fmt2 p2: fmt p1", err.Error())
	assert.True(t, err.Contains("fmt2 %v"))
	assert.False(t, err.Contains("fmt2 p2"))
	assert.True(t, err.Contains("fmt %v"))
	assert.False(t, err.Contains("fmt p1"))
}

func TestContains_TopLevelNilErr(t *testing.T) {
	assert.False(t, xerror.Contains(nil, "msg"))
}

func TestContains_TopLevelNativeErr(t *testing.T) {
	err := errors.New("msg")
	assert.True(t, xerror.Contains(err, "msg"))
	assert.False(t, xerror.Contains(err, "else"))
}

func TestContains_TopLevelError(t *testing.T) {
	err := xerror.Wrap(xerror.New("fmt %v", "p1"), "fmt2 %v", "p2")
	assert.Equal(t, "fmt2 p2: fmt p1", err.Error())
	assert.True(t, xerror.Contains(err, "fmt2 %v"))
	assert.False(t, xerror.Contains(err, "fmt2 p2"))
	assert.True(t, xerror.Contains(err, "fmt %v"))
	assert.False(t, xerror.Contains(err, "fmt p1"))
}

func TestClone_FormatOnly(t *testing.T) {
	err := xerror.New("fmt")
	cp := err.Clone()
	assert.Equal(t, err, cp)
	assert.True(t, err != cp)
}

func TestClone_PlaceholdersAndDebug(t *testing.T) {
	err := xerror.New("fmt %v %v", "p2", "p1", "d2", "d1")
	cp := err.Clone()
	assert.Equal(t, err, cp)
	assert.True(t, err != cp)
}

func TestImplementsError(t *testing.T) {
	var err error
	err = xerror.New("fmt")
	assert.Equal(t, "fmt", err.Error())
	_, ok := err.(xerror.Error)
	assert.True(t, ok)
}

func TestMarshalJSON(t *testing.T) {
	_, err := json.Marshal(xerror.New("fmt %v", "p1", "d2", "d1"))
	assert.Nil(t, err)
}
