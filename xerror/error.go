package xerror

import "strings"

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

// Wrap creates an augmented error given a standard Go error.
func Wrap(err error) *Error {
	return New(err.Error())
}

// Is returns true if the outermost error message equals the given message, false otherwise.
func (e *Error) Is(message string) bool {
	return len(e.messages) > 0 && e.messages[0] == message
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
