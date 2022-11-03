package posixmq

// Queue models a message queue.
type Queue interface {
	Receive() (data []byte, err error)
	Send(data []byte, priority uint) error
	Unlink(qname string) error
}
