// Implement something close to what the very cool package
//
//   github.com/pkg/errors
//
// does, but also support embedding HTTP error codes in the
// returned errors.
package errors

import (
	goerrors "errors"
	"fmt"
	"strings"
)

type httperror struct {
	text   string
	status int
}

// Return the description of this error. Implements
// the error interface.
func (he *httperror) Error() string {
	return he.text
}

// Create a new error with text and no HTTP status code
// set.
func New(text string) error {
	tt := strings.TrimSpace(text)
	return goerrors.New(tt)
}

// Create a new error with format string and no HTTP status
// code.
func Newf(format string, formatargs ...interface{}) error {
	text := fmt.Sprintf(format, formatargs...)
	return New(text)
}

// Create a new error with text and HTTP status code set.
func NewWith(status int, text string) error {
	tt := strings.TrimSpace(text)

	return &httperror{
		text:   tt,
		status: status,
	}
}

// Create a new error with format string and an HTTP status
// code set.
func NewfWith(status int, format string, formatargs ...interface{}) error {
	err := Newf(format, formatargs...)
	return WithStatus(status, err)
}

// Wrap cause and annotate it with text.
func Wrap(cause error, text string) error {
	w := fmt.Sprintf("%s: %s", text, cause.Error())

	// preserve http status code if there is one set
	if status, ok := Status(cause); ok {
		return NewWith(status, w)
	} else {
		return New(w)
	}
}

// Wrap cause and annotate it with text and HTTP status.
func WrapWith(status int, cause error, text string) error {
	err := Wrap(cause, text)
	return WithStatus(status, err)
}

// Wrap cause and annotate it with a format string.
func Wrapf(cause error, format string, formatargs ...interface{}) error {
	text := fmt.Sprintf(format, formatargs...)
	return Wrap(cause, text)
}

// Wrap cause and annotate it with an HTTP status and a format
// string.
func WrapfWith(status int, cause error, format string, formatargs ...interface{}) error {
	err := Wrapf(cause, format, formatargs...)
	return WithStatus(status, err)
}

// If err does in fact wrap an HTTP status code, return
// that status code.
func Status(err error) (status int, ok bool) {
	if he, ok := err.(*httperror); ok {
		return he.status, true
	} else {
		return 0, false
	}
}

// Return err with HTTP status attached. If err already carried
// an HTTP status, it is overwritten.
func WithStatus(status int, err error) error {
	return NewWith(status, err.Error())
}
