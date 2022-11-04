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

const (
	mq_open    = unix.SYS_MQ_OPEN
	mq_send    = unix.SYS_MQ_TIMEDSEND
	mq_receive = unix.SYS_MQ_TIMEDRECEIVE
	mq_unlink  = unix.SYS_MQ_UNLINK
)

// MessageQueue
type MessageQueue struct {
	name string //name of the queue
	fd   int    //file descriptor for the queue
}

type mqAttr struct {
	_              int64
	MaxQueueSize   int64
	MaxMessageSize int64
	_              int64
}

// Open creates a message queue for read and write.
func Open(qname string) (m *MessageQueue, err error) {

	m = new(MessageQueue)
	//TODO: this is weird, because according to spec names should start with /
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
	m.fd = int(mqd)
	return
}
