/*
Package xerror extends the functionality of Go's built-in error interface, in several ways:

	- errors carry a stack trace of the location where they were first created
	- it is possible to attach debug objects to errors for later reporting
	- errors hold a list of messages, and can be wrapped by prepending new messages to the list
	- errors can natively be treated as Go built-in errors and serialized to a string message or JSON representation

To create a new error, use the New function:

	err := xerror.New("first message", "second message")

When this error is converted to string using the Error method from Go's error interface, the following is returned:

	"first message: second message"

To create an augmented error given a Go error, use the Wrap function:

	if _, err := pkg.Method(...); err != nil {
		return xerror.Wrap(err).WithMessages("unable to execute Method")
	}

If the given error is actually of type Error, the Error is immediately returned unmodified.

Errors are immutable, but modified copies of them can be obtained using WithMessages and WithDebug:

	return xerror.Wrap(err).WithMessages("unable to perform action").WithDebug(ctx, req)

Error also provides methods for determining its type, by matching all messages or just the outermost one:

	err := xerror.New("m2", "m1")
	err.Is("m2") // true
	err.Is("m1") // false
	err.IsPattern(regexp.MustCompile("2$")) // true
	err.IsPattern(regexp.MustCompile("1$")) // false
	err.Contains("m2") // true
	err.Contains("m1") // true
	err.ContainsPattern(regexp.MustCompile("2$") // true
	err.ContainsPattern(regexp.MustCompile("1$") // true
*/
package xerror

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Error is the augmented error interface provided by this package.
type Error interface {
	error
	json.Marshaler

	Is(string) bool
	IsPattern(*regexp.Regexp) bool
	Contains(string) bool
	ContainsPattern(*regexp.Regexp) bool
	Debug() []interface{}
	Stack() []string
	Clone() Error
}

// xerror is the internal implementation of Error
type xerror struct {
	msg   string
	fmts  []string
	dbg   []interface{}
	stack []string
}

// xerrorJSON is used to serialize Error to JSON
type xerrorJSON struct {
	Message string        `json:"message"`
	Debug   []interface{} `json:"debug,omitempty"`
	Stack   []string      `json:"stack"`
}

// New returns a new augmented error. Parameters that don't have a placeholder in the format string are only stored as debug objects.
func New(format string, v ...interface{}) Error {
	v = nilToEmpty(v)
	return &xerror{
		msg:   safeSprintf(format, v),
		fmts:  []string{format},
		dbg:   v,
		stack: newStack(),
	}
}

// Wrap returns a new augmented error that wraps the given Go `error` or `Error`.
func Wrap(err error, format string, v ...interface{}) Error {
	v = nilToEmpty(v)
	xerr := cloneOrNew(err)
	xerr.msg = fmt.Sprintf("%v: %v", safeSprintf(format, v), xerr.msg)
	xerr.fmts = append([]string{format}, xerr.fmts...)
	xerr.dbg = append(v, xerr.dbg...)
	return xerr
}

// Error implements the `error` interface.
func (e *xerror) Error() string {
	return e.msg
}

// MarshalJSON implements the `json.Marshaler` interface.
func (e *xerror) MarshalJSON() ([]byte, error) {
	return json.Marshal(&xerrorJSON{
		Message: e.msg,
		Debug:   e.dbg,
		Stack:   e.stack,
	})
}

// Is returns true if the outermost error message format equals the given message format, false otherwise.
func (e *xerror) Is(fmt string) bool {
	return e.fmts[0] == fmt
}

// IsPattern returns true if the outermost error message format matches the given pattern, false otherwise.
func (e *xerror) IsPattern(pattern *regexp.Regexp) bool {
	return pattern.MatchString(e.fmts[0])
}

// Contains returns true if the error contains the given message format, false otherwise.
func (e *xerror) Contains(format string) bool {
	for _, f := range e.fmts {
		if f == format {
			return true
		}
	}
	return false
}

// ContainsPattern returns true if the error contains a message format that matches the given pattern, false otherwise.
func (e *xerror) ContainsPattern(pattern *regexp.Regexp) bool {
	for _, f := range e.fmts {
		if pattern.MatchString(f) {
			return true
		}
	}
	return false
}

// Debug returns the slice of debug objects.
func (e *xerror) Debug() []interface{} {
	return e.dbg
}

// Stack returns the stack trace associated with the error.
func (e *xerror) Stack() []string {
	return e.stack
}

// Clone returns an exact copy of the `Error`.
func (e *xerror) Clone() Error {
	return &xerror{
		msg:   e.msg,
		fmts:  append(make([]string, 0, len(e.fmts)), e.fmts...),
		dbg:   append(make([]interface{}, 0, len(e.dbg)), e.dbg...),
		stack: append(make([]string, 0, len(e.stack)), e.stack...),
	}
}

// Is returns true if the outermost message format (if `err` is `Error`) or error string (if `err` is a Go `error`) equals the given message.
func Is(err error, message string) bool {
	if xerr, ok := err.(*xerror); ok {
		return xerr.Is(message)
	}
	return err.Error() == message
}

// IsPattern is like Is but uses regexp matching rather than string comparison.
func IsPattern(err error, pattern *regexp.Regexp) bool {
	if xerr, ok := err.(*xerror); ok {
		return xerr.IsPattern(pattern)
	}
	return pattern.MatchString(err.Error())
}

// Contains is like Is, but in case `err` is of type `Error` compares the message format with all attached message formats.
func Contains(err error, message string) bool {
	if xerr, ok := err.(*xerror); ok {
		return xerr.Contains(message)
	}
	return err.Error() == message
}

// ContainsPattern is like Contains but uses regexp matching rather than string comparison.
func ContainsPattern(err error, pattern *regexp.Regexp) bool {
	if xerr, ok := err.(*xerror); ok {
		return xerr.ContainsPattern(pattern)
	}
	return pattern.MatchString(err.Error())
}

// cloneOrNew wraps the given `error` unless it is already of type `*xerror`, in which case it returns a copy
func cloneOrNew(err error) *xerror {
	if xerr, ok := err.(*xerror); ok {
		return xerr.Clone().(*xerror)
	}
	return New(err.Error()).(*xerror)
}

// safeSprintf is like `fmt.Sprintf`, but passes through only at most parameters as placeholders in the format string
func safeSprintf(format string, v []interface{}) string {
	if n := strings.Count(format, "%") - strings.Count(format, "%%")*2; len(v) > n {
		v = v[:n]
	}
	return fmt.Sprintf(format, v...)
}

// nilToEmpty returns the given slice if not nil, or an empty slice if nil
func nilToEmpty(v []interface{}) []interface{} {
	if v == nil {
		return []interface{}{}
	}
	return v
}
