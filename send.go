//go:build linux
// +build linux

package posixmq

import (
	"errors"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Ensure type implements the interface.
var _ Queue = (*MessageQueue)(nil)

// MessageQueueOpenMode is the mode used to open message queues.
const MessageQueueOpenMode = 0640

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

// Open creates a message queue for read and write.
func Open(qname string) (m *MessageQueue, err error) {

	m = &MessageQueue{}
	//TODO: this is weird, because according to spec names should start with /
	if strings.HasPrefix(qname, "/") {
		return m, errors.New("message queue must start with /")
	}

	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return nil, err
	}

	flags := unix.O_RDWR | unix.O_CREAT

	//mode := unix.S_IRUSR | unix.S_IWUSR | unix.S_IRGRP | unix.S_IWGRP | unix.S_IROTH | unix.S_IWOTH

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
		})), //queue attributes
		0, //unused
		0, //unused
	)
	if errno != 0 {
		return nil, errno
	}
	m.name = qname
	m.fd = mqd
	return
}

func (m *MessageQueue) Receive() ([]byte, error) {

	//mqd_t mqdes, char *restrict msg_ptr,
	//size_t msg_len, unsigned int *restrict msg_prio,
	//const struct timespec *restrict abs_timeout

	deadline := time.Now().Add(3 * time.Second)

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
		uintptr(0),                        // msg_prio
		uintptr(unsafe.Pointer(&t)),       // abs_timeout
		0,                                 // unused
	)
	if errno != 0 {
		return nil, errno
	}
	return data, nil

}

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
