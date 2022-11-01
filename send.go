//go:build linux
// +build linux

package posixmq

import (
	"context"
	"errors"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

// MessageQueueOpenMode is the mode used to open message queues.
const MessageQueueOpenMode = 0644

type messageQueue uintptr

type mqAttr struct {
	_              int64
	MaxQueueSize   int64
	MaxMessageSize int64
	_              int64
}

func Open(qname string) (m MessageQueue, err error) {

	if strings.HasPrefix(qname, "/") {
		return m, errors.New("queue shall not start with /")
	}

	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return nil, err
	}

	flags := unix.O_WRONLY | unix.O_CREAT

	// From MQ_OPEN(3) manpage:
	// mqd_t mq_open(const char *name, int oflag, mode_t mode, struct mq_attr *attr);
	mqd, _, errno := unix.Syscall6(
		unix.SYS_MQ_OPEN,
		uintptr(unsafe.Pointer(name)), // name
		uintptr(flags),                // oflag
		uintptr(MessageQueueOpenMode), // mode
		uintptr(unsafe.Pointer(&mqAttr{
			MaxQueueSize:   10,
			MaxMessageSize: 8192,
		})), 0, 0,
	)
	if errno != 0 {
		return nil, errno
	}
	return messageQueue(mqd), nil
}

func (m messageQueue) Receive(ctx context.Context) (data []byte, err error) {

	return
}

func (m messageQueue) Send(ctx context.Context, data []byte) (err error) {

	// From MQ_SEND(3) manpage:
	// mqd_t mqdes, const char *msg_ptr size_t msg_len, unsigned int msg_prio
	mqd, _, errno := unix.Syscall6(
		unix.SYS_MQ_OPEN,
		uintptr(unsafe.Pointer(name)), // name
		uintptr(flags),                // oflag
		uintptr(MessageQueueOpenMode), // mode
		uintptr(unsafe.Pointer(&mqAttr{
			MaxQueueSize:   10,
			MaxMessageSize: 8192,
		})), 0, 0,
	)
	if errno != 0 {
		return nil, errno
	}
	return
}
