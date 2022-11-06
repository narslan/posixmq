package posixmq

import (
	"context"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func (m *MessageQueue) Receive(ctx context.Context) ([]byte, error) {

	deadline := time.Now().Add(-1)

	t, err := unix.TimeToTimespec(deadline)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 8192)
	_, _, errno := unix.Syscall6(
		mq_receive,
		uintptr(m.fd),                     // mqdes
		uintptr(unsafe.Pointer(&data[0])), // msg_ptr
		uintptr(8192),                     // msg_len
		0,                                 // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)
	if errno != 0 {
		return nil, errno
	}
	return data, nil

}
