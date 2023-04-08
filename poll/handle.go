package poll

import (
	"os"

	"github.com/narslan/posixmq"
)

// Desc is a file access descriptor.
// It's methods are not goroutine safe.
type Desc struct {
	file  int
	event Event
}

func (d *Desc) fd() int {
	return d.file
}

// Must is a helper that wraps a call to a function returning (*Desc, error).
// It panics if the error is non-nil and returns desc if not.
// It is intended for use in short Desc initializations.
func Must(desc *Desc, err error) *Desc {
	if err != nil {
		panic(err)
	}
	return desc
}

// HandleRead creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventEdgeTriggered).
func HandleRead(conn *posixmq.MessageQueue) (*Desc, error) {
	return Handle(conn, EventRead|EventEdgeTriggered)
}

// HandleReadOnce creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventOneShot).
func HandleReadOnce(conn *posixmq.MessageQueue) (*Desc, error) {
	return Handle(conn, EventRead|EventOneShot)
}

// HandleWrite creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventEdgeTriggered).
func HandleWrite(conn *posixmq.MessageQueue) (*Desc, error) {
	return Handle(conn, EventWrite|EventEdgeTriggered)
}

// HandleWriteOnce creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventOneShot).
func HandleWriteOnce(conn *posixmq.MessageQueue) (*Desc, error) {
	return Handle(conn, EventWrite|EventOneShot)
}

// HandleReadWrite creates read and write descriptor for further use in Poller
// methods.
// It is the same as Handle(conn, EventRead|EventWrite|EventEdgeTriggered).
func HandleReadWrite(conn *posixmq.MessageQueue) (*Desc, error) {
	return Handle(conn, EventRead|EventWrite|EventEdgeTriggered)
}

// Handle creates new Desc with given conn and event.
// Returned descriptor could be used as argument to Start(), Resume() and
// Stop() methods of some Poller implementation.
func Handle(mq *posixmq.MessageQueue, event Event) (*Desc, error) {
	desc, err := handle(mq, event)
	if err != nil {
		return nil, err
	}

	if err = setNonblock(desc.fd(), true); err != nil {
		return nil, os.NewSyscallError("setnonblock", err)
	}

	return desc, nil
}

func handle(mq *posixmq.MessageQueue, event Event) (*Desc, error) {
	// f, ok := x.(filer)
	// if !ok {
	// 	return nil, ErrNotFiler
	// }

	// // Get a copy of fd.
	// file, err := f.File()
	// if err != nil {
	// 	return nil, err
	// }

	return &Desc{
		file:  mq.FD,
		event: event,
	}, nil
}
