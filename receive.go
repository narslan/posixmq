package posixmq

import (
	"bytes"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func (m *MessageQueue) Receive() ([]byte, error) {
	deadline := time.Now().Add(1 * time.Minute)

	t, err := unix.TimeToTimespec(deadline)
	if err != nil {
		return nil, err
	}

	data := make([]byte, m.MessageSize)
	_, _, errno := unix.Syscall6(
		mqReceive,
		uintptr(m.FD),                     // mqdes
		uintptr(unsafe.Pointer(&data[0])), // msg_ptr
		uintptr(m.MessageSize),            // msg_len
		0,                                 // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)

	switch errno {
	default:
		return nil, errno
	case 0:
		return bytes.TrimRight(data, "\x00"), nil
	case syscall.EMSGSIZE:
		return nil, MessageTooLongError{
			MessageSize: len(data),
			Cause:       errno,
		}
	case syscall.ETIMEDOUT:
		return nil, TimedOutError{
			Cause: errno,
		}
	}
}
