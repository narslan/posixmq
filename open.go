package posixmq

import (
	"context"
	"fmt"
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
	*Config
	//The file descriptor for the queue
	fd int
}

// mqAttr is the attributes of message queue.
type mqAttr struct {
	_              int64
	MaxQueueSize   int64
	MaxMessageSize int64
	_              int64
}

type Config struct {
	QueueSize   int64
	MessageSize int64
	Name        string
}

// Open creates a message queue for read and write.
func Open(ctx context.Context, cfg *Config) (m *MessageQueue, err error) {

	m = new(MessageQueue)
	m.Config = cfg

	//BytePtrFromString returns a pointer to a NUL-terminated array of bytes
	name, err := unix.BytePtrFromString(cfg.Name)
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
			MaxQueueSize:   cfg.QueueSize,
			MaxMessageSize: cfg.MessageSize,
		})), //queue attributes
		0, //unused
		0, //unused
	)
	switch errno {
	case 0:
		m.fd = int(mqd)
		return m, nil
	default:
		return nil, fmt.Errorf("[open] %w", errno)
	}

}

// Close closes the message que queue.
func (m *MessageQueue) Close(ctx context.Context) error {
	return unix.Close(int(m.fd))
}

func (m *MessageQueue) Unlink(ctx context.Context, qname string) error {

	name, err := unix.BytePtrFromString(qname)
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
