//go:build linux
// +build linux

package posixmq

import (
	"context"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func (m *MessageQueue) Send(ctx context.Context, data []byte, priority uint) (err error) {
	// From MQ_SEND(3) manpage, regarding mq_timedsend:
	//
	//     If the message queue is full, and the timeout has already expired by
	//     the time of the call, mq_timedsend() returns immediately.
	//
	// Setting a timeout to the past to enable non-blocking mode.
	deadline := time.Now().Add(-1)

	t, err := unix.TimeToTimespec(deadline)
	if err != nil {
		return err
	}

	// From MQ_SEND(3) manpage:
	// mqd_t mqdes, const char *msg_ptr size_t msg_len, unsigned int msg_prio
	_, _, errno := unix.Syscall6(
		mqSend,
		uintptr(m.FD),                     // mqdes
		uintptr(unsafe.Pointer(&data[0])), // msg_ptr
		uintptr(len(data)),                // msg_len
		uintptr(priority),                 // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)
	switch errno {
	default:
		return errno
	case 0:
		return nil
	case syscall.EMSGSIZE:
		return MessageTooLongError{
			MessageSize: len(data),
			Cause:       errno,
		}
	case syscall.EAGAIN:
		return fmt.Errorf("again: %s", errno)
	case syscall.ETIMEDOUT:
		return TimedOutError{
			Cause: errno,
		}

	}
}
