//go:build linux
// +build linux

package posixmq

import (
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func (m *MessageQueue) Send(data []byte, priority uint) (err error) {

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
		mq_send,
		uintptr(m.fd),                     // mqdes
		uintptr(unsafe.Pointer(&data[0])), // msg_ptr
		uintptr(len(data)),                // msg_len
		uintptr(priority),                 // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)
	if errno != 0 {
		return
	}
	return
}

func (m *MessageQueue) Close() error {
	return unix.Close(int(m.fd))
}

func (m *MessageQueue) Unlink(qname string) error {

	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return err
	}
	// Close via the file descriptor before removing the queue.
	err = m.Close()
	if err != nil {
		return err
	}

	_, _, errno := unix.Syscall(
		mq_unlink,
		uintptr(unsafe.Pointer(name)), // name
		0,
		0,
	)

	if errno != 0 {
		return errno
	}
	return nil
}
