package posixmq

import (
	"fmt"
	"strings"
)

//TimedOutError is the error returned by mqsend because of the queue was full.

type TimedOutError struct {
	Cause error
}

func (e TimedOutError) Error() string {
	return fmt.Sprintf("mqsend: send timed out: %v", e.Cause)
}

// Unwrap returns the underlying error.
func (e TimedOutError) Unwrap() error {
	return e.Cause
}

// MessageTooLargeError is the error returned by MessageQueue.Send when the
// message is larger than the configured max size.

type MessageTooLongError struct {
	MessageSize int

	// Note that on non-linux systems Cause will be nil.
	Cause error
}

func (e MessageTooLongError) Error() string {
	var sb strings.Builder

	sb.WriteString(e.Cause.Error())

	return sb.String()
}

// Unwrap returns the underlying error, if any.
func (e MessageTooLongError) Unwrap() error {
	return e.Cause
}
