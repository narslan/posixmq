package posixmq

import "context"

// Queue models a message queue.
type Queue interface {
	Receive(ctx context.Context) (data []byte, err error)
	Send(ctx context.Context, data []byte) error
}
