package posixmq

import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Ensure type implements the interface.
var _ Queue = (*MessageQueue)(nil)

// MessageQueueOpenMode is the mode used to open message queues.
const MessageQueueOpenMode = 0o640

// Aliasing signals.
const (
	mqOpen    = unix.SYS_MQ_OPEN
	mqSend    = unix.SYS_MQ_TIMEDSEND
	mqReceive = unix.SYS_MQ_TIMEDRECEIVE
	mqUnlink  = unix.SYS_MQ_UNLINK
)

// MessageQueue embraces attributes of a message queue.
type MessageQueue struct {
	*Config
	// The file descriptor for the queue
	FD int
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

func (cfg *Config) String() string {
	var b strings.Builder

	b.WriteString("queue message size queue size \n")

	fmt.Fprintf(&b, "queue %q has [message size]:%d [queue capacity]:%d \n", cfg.Name, cfg.MessageSize, cfg.QueueSize)
	return b.String()
}

// Open creates a message queue for read and write.
func Open(ctx context.Context, cfg *Config) (m *MessageQueue, err error) {
	m = new(MessageQueue)
	m.Config = cfg

	// BytePtrFromString returns a pointer to a NUL-terminated array of bytes
	name, err := unix.BytePtrFromString(cfg.Name)
	if err != nil {
		return nil, err
	}

	flags := unix.O_RDWR | unix.O_CREAT

	// From MQ_OPEN(3) manpage:
	// mqd_t mq_open(const char *name, int oflag, mode_t mode, struct mq_attr *attr);
	mqd, _, errno := unix.Syscall6(
		mqOpen,
		uintptr(unsafe.Pointer(name)), // name
		uintptr(flags),                // oflag
		uintptr(MessageQueueOpenMode), // mode
		uintptr(unsafe.Pointer(&mqAttr{
			MaxQueueSize:   m.QueueSize,
			MaxMessageSize: m.MessageSize,
		})), // queue attributes
		0, // unused
		0, // unused
	)
	switch errno {
	case 0:
		m.FD = int(mqd)
		return
	default:
		return nil, fmt.Errorf("[open] %w", errno)

	}
}

// Close closes the message queue.
func (m *MessageQueue) Close(ctx context.Context) error {
	return unix.Close(int(m.FD))
}

func (m *MessageQueue) Unlink(ctx context.Context, qname string) error {
	name, err := unix.BytePtrFromString(qname)
	if err != nil {
		return err
	}

	_, _, errno := unix.Syscall(
		mqUnlink,
		uintptr(unsafe.Pointer(name)), // name
		0,
		0,
	)

	if errno != 0 {
		return errno
	}
	return nil
}
