//go:build linux
// +build linux

package posixmq

import (
	"context"
	"errors"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// MessageQueueOpenMode is the mode used to open message queues.
const MessageQueueOpenMode = 0644

const (
	mq_open    = unix.SYS_MQ_OPEN
	mq_send    = unix.SYS_MQ_TIMEDSEND
	mq_receive = unix.SYS_MQ_TIMEDRECEIVE
	mq_unlink  = unix.SYS_MQ_UNLINK
)

// MessageQueue
type MessageQueue struct {
	name string  //name of the queue
	fd   uintptr //file descriptor for the queue
}

type mqAttr struct {
	_              int64
	MaxQueueSize   int64
	MaxMessageSize int64
	_              int64
}

func Open(qname string) (m *MessageQueue, err error) {

	m = &MessageQueue{}
	if strings.HasPrefix(qname, "/") {
		return m, errors.New("message queue shall not start with /")
	}

	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return nil, err
	}

	flags := unix.O_WRONLY | unix.O_CREAT

	// From MQ_OPEN(3) manpage:
	// mqd_t mq_open(const char *name, int oflag, mode_t mode, struct mq_attr *attr);
	mqd, _, errno := unix.Syscall6(
		mq_open,
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
	m.name = qname
	m.fd = mqd
	return
}

func (m *MessageQueue) Receive(ctx context.Context) (data []byte, err error) {

	return
}

func (m *MessageQueue) Send(ctx context.Context, data []byte) (err error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		// From MQ_SEND(3) manpage, regarding mq_timedsend:
		//
		//     If the message queue is full, and the timeout has already expired by
		//     the time of the call, mq_timedsend() returns immediately.
		//
		// So set a timeout to the past to enable non-blocking mode when there's no
		// deadline set in the context object.
		deadline = time.Now().Add(-1)
	}

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
		0,                                 // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)
	if errno != 0 {
		return
	}
	return
}
