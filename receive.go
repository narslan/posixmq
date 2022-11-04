package posixmq

import (
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func (m *MessageQueue) Receive() ([]byte, error) {

	//From MQ_TIMEDRECEIVE(3) manpage, regarding mq_timedreceive:
	//mqd_t mqdes, char *restrict msg_ptr,
	//size_t msg_len, unsigned int *restrict msg_prio,
	//const struct timespec *restrict abs_timeout

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
