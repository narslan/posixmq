package posixmq

// Queue is a model for a message queue.
type Queue interface {
	Receive() (data []byte, err error)
	Send(data []byte, priority uint) error
	Unlink() error
}
