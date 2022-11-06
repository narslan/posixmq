package posixmq

import "context"

// Queue is a model for a message queue.
type Queue interface {
	Receive(ctx context.Context) (data []byte, err error)
	Send(ctx context.Context, data []byte, priority uint) error
	Unlink(ctx context.Context, qname string) error
}
