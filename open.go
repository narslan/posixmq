package posixmq

import (
	"errors"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Ensure type implements the interface.
var _ Queue = (*MessageQueue)(nil)

// MessageQueueOpenMode is the mode used to open message queues.
const MessageQueueOpenMode = 0640

// Aliasing signals.
const (
	mq_open    = unix.SYS_MQ_OPEN
	mq_send    = unix.SYS_MQ_TIMEDSEND
	mq_receive = unix.SYS_MQ_TIMEDRECEIVE
	mq_unlink  = unix.SYS_MQ_UNLINK
)

// MessageQueue embraces attributes of a message queue.
type MessageQueue struct {
	//Name of the queue
	name string

	//The file descriptor for the queue
	fd int

	// The max number of messages in the queue.
	queueSize int64

	// The max size in bytes per message.
	messageSize int64
}

type mqAttr struct {
	_              int64
	MaxQueueSize   int64
	MaxMessageSize int64
	_              int64
}

// Open creates a message queue for read and write.
func Open(qname string, options ...func(*MessageQueue)) (m *MessageQueue, err error) {

	m = new(MessageQueue)

	m.queueSize = 10
	m.messageSize = 8192

	for _, applyOpt := range options {
		applyOpt(m)
	}

	//TODO: this is weird, because according to the spec names should start with /
	if strings.HasPrefix(qname, "/") {
		return m, errors.New("message queue must not start with /")
	}

	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return nil, err
	}

	flags := unix.O_RDWR | unix.O_CREAT

	// From MQ_OPEN(3) manpage:
	// mqd_t mq_open(const char *name, int oflag, mode_t mode, struct mq_attr *attr);
	mqd, _, errno := unix.Syscall6(
		mq_open,
		uintptr(unsafe.Pointer(name)), // name
		uintptr(flags),                // oflag
		uintptr(MessageQueueOpenMode), // mode
		uintptr(unsafe.Pointer(&mqAttr{
			MaxQueueSize:   m.queueSize,
			MaxMessageSize: m.messageSize,
		})), //queue attributes
		0, //unused
		0, //unused
	)
	if errno != 0 {
		return nil, errno
	}
	m.name = qname
	m.fd = int(mqd)
	return
}

func (m *MessageQueue) WithQueueSize(n int64) func(*MessageQueue) {
	return func(m *MessageQueue) {
		m.queueSize = n
	}
}

func (m *MessageQueue) WithMessageSize(n int64) func(*MessageQueue) {
	return func(m *MessageQueue) {
		m.messageSize = n
	}
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
