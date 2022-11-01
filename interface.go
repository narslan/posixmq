package posixmq

import "context"

// MessageQueue describes a posix mq.
type MessageQueue interface {
	Receive(ctx context.Context) (data []byte, err error)
	Send(ctx context.Context, data []byte) error
}
