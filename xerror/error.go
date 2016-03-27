/*
Package xerror extends the functionality of Go's built-in error interface.


*/
package xerror

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Error represents an augmented error.
type Error interface {
	error
	json.Marshaler

	Is(string) bool
	IsPattern(*regexp.Regexp) bool
	Contains(string) bool
	ContainsPattern(*regexp.Regexp) bool
	Messages() []string
	Debug() []interface{}
	Stack() []string
	Copy() Error
	WithMessages(...string) Error
	WithDebug(...interface{}) Error
}

// xerror implements Error
type xerror struct {
	messages []string
	debug    []interface{}
	stack    []string
}

// xerrorJSON is used to produce the JSON representation of an Error
type xerrorJSON struct {
	Messages []string      `json:"messages"`
	Debug    []interface{} `json:"debug,omitempty"`
	Stack    []string      `json:"stack"`
}

// New creates an augmented error given a list of messages.
func New(messages ...string) Error {
	return &xerror{
		messages: messages,
		debug:    make([]interface{}, 0),
		stack:    newStack(),
	}
}

// Wrap creates an augmented error given a standard Go error or just returns the given *Error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(*xerror); ok {
		return xerr
	}
	return New(err.Error())
}

// Error implements the standard error interface.
// The result is built by joining the messages with the ": " separator.
func (e *xerror) Error() string {
	return strings.Join(e.messages, ": ")
}

// MarshalJSON implements the JSON Marshaler interface.
func (e *xerror) MarshalJSON() ([]byte, error) {
	return json.Marshal(&xerrorJSON{
		Messages: e.messages,
		Debug:    e.debug,
		Stack:    e.stack,
	})
}

// Is returns true if the outermost error message equals the given message, false otherwise.
func (e *xerror) Is(message string) bool {
	return e.messages[0] == message
}

// IsPattern returns true if the outermost error message matches the given pattern, false otherwise.
func (e *xerror) IsPattern(pattern *regexp.Regexp) bool {
	return pattern.MatchString(e.messages[0])
}

// Contains returns true if the error contains the given message, false otherwise.
func (e *xerror) Contains(message string) bool {
	for _, m := range e.messages {
		if m == message {
			return true
		}
	}
	return false
}

// ContainsPattern returns true if the error contains a message that matches the given pattern, false otherwise.
func (e *xerror) ContainsPattern(pattern *regexp.Regexp) bool {
	for _, m := range e.messages {
		if pattern.MatchString(m) {
			return true
		}
	}
	return false
}

// Messages returns the slice of error messages.
func (e *xerror) Messages() []string {
	return e.messages
}

// Debug returns the slice of debug objects.
func (e *xerror) Debug() []interface{} {
	return e.debug
}

// Stack returns the innermost error stack trace.
func (e *xerror) Stack() []string {
	return e.stack
}

// Copy returns a copy of the error.
func (e *xerror) Copy() Error {
	return &xerror{
		messages: append(make([]string, 0, len(e.messages)), e.messages...),
		debug:    append(make([]interface{}, 0, len(e.debug)), e.debug...),
		stack:    append(make([]string, 0, len(e.stack)), e.stack...),
	}
}

// WithMessages returns a copy of the Error with the given messages prepended to the messages slice.
func (e *xerror) WithMessages(message ...string) Error {
	n := e.Copy().(*xerror)
	n.messages = append(message, n.messages...)
	return n
}

// WithDebug returns a copy of the Error with the given debug objects prepended to the debug objects slice.
func (e *xerror) WithDebug(debug ...interface{}) Error {
	n := e.Copy().(*xerror)
	n.debug = append(debug, n.debug...)
	return n
}

// Is returns true if the outermost error message (if err is *Error) or the error string (if err is a standard Go error) equals the given message.
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

// Contains is like Is, but in case err is of type *Error compares the message with all attached messages.
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
