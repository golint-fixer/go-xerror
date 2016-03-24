/*
Package xerror extends the functionality of Go's built-in error interface in few ways:

1) A list of error messages, for easy wrapping and error type detection
2) Ability to attach arbitrary debug values to the error
3) Error stack trace propagation
4) Immutable errors that can easily be copied and modified with a fluent API
*/
package xerror

import (
	"regexp"
	"strings"
)

// Error represents an augmented error.
type Error struct {
	messages []string
	debug    []interface{}
	stack    []string
}

// New creates an augmented error given a list of messages.
func New(messages ...string) *Error {
	return &Error{
		messages: messages,
		debug:    make([]interface{}, 0),
		stack:    newStack(),
	}
}

// Wrap creates an augmented error given a standard Go error or just returns the given *Error.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(*Error); ok {
		return xerr
	}
	return New(err.Error())
}

// Is returns true if the outermost error message equals the given message, false otherwise.
func (e *Error) Is(message string) bool {
	return e.messages[0] == message
}

// IsPattern returns true if the outermost error message matches the given pattern, false otherwise.
func (e *Error) IsPattern(pattern *regexp.Regexp) bool {
	return pattern.MatchString(e.messages[0])
}

// Contains returns true if the error contains the given message, false otherwise.
func (e *Error) Contains(message string) bool {
	for _, m := range e.messages {
		if m == message {
			return true
		}
	}
	return false
}

// ContainsPattern returns true if the error contains a message that matches the given pattern, false otherwise.
func (e *Error) ContainsPattern(pattern *regexp.Regexp) bool {
	for _, m := range e.messages {
		if pattern.MatchString(m) {
			return true
		}
	}
	return false
}

// Error implements the standard error interface.
// The result is built by joining the messages with the ": " separator.
func (e *Error) Error() string {
	return strings.Join(e.messages, ": ")
}

// Messages returns the slice of error messages.
func (e *Error) Messages() []string {
	return e.messages
}

// Debug returns the slice of debug objects.
func (e *Error) Debug() []interface{} {
	return e.debug
}

// Stack returns the innermost error stack trace.
func (e *Error) Stack() []string {
	return e.stack
}

// Copy returns a copy of the error.
func (e *Error) Copy() *Error {
	return &Error{
		messages: append(make([]string, 0, len(e.messages)), e.messages...),
		debug:    append(make([]interface{}, 0, len(e.debug)), e.debug...),
		stack:    append(make([]string, 0, len(e.stack)), e.stack...),
	}
}

// WithMessages returns a copy of the Error with the given messages prepended to the messages slice.
func (e *Error) WithMessages(message ...string) *Error {
	n := e.Copy()
	n.messages = append(message, n.messages...)
	return n
}

// WithDebug returns a copy of the Error with the given debug objects prepended to the debug objects slice.
func (e *Error) WithDebug(debug ...interface{}) *Error {
	n := e.Copy()
	n.debug = append(debug, n.debug...)
	return n
}

// Is returns true if the outermost error message (if err is *Error) or the error string (if err is a standard Go error) equals the given message.
func Is(err error, message string) bool {
	if xerr, ok := err.(*Error); ok {
		return xerr.Is(message)
	} else {
		return err.Error() == message
	}
}

// IsPattern is like Is but uses regexp matching rather than string comparison.
func IsPattern(err error, pattern *regexp.Regexp) bool {
	if xerr, ok := err.(*Error); ok {
		return xerr.IsPattern(pattern)
	} else {
		return pattern.MatchString(err.Error())
	}
}

// Contains is like Is, but in case err is of type *Error compares the message with all attached messages.
func Contains(err error, message string) bool {
	if xerr, ok := err.(*Error); ok {
		return xerr.Contains(message)
	} else {
		return err.Error() == message
	}
}

// ContainsPattern is like Contains but uses regexp matching rather than string comparison.
func ContainsPattern(err error, pattern *regexp.Regexp) bool {
	if xerr, ok := err.(*Error); ok {
		return xerr.ContainsPattern(pattern)
	} else {
		return pattern.MatchString(err.Error())
	}
}
