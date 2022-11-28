package posixmq_test

import (
	"context"
	"github.com/narslan/posixmq"
	"testing"
)

func TestOpenWithEpoll(t *testing.T) {

	ctx := context.Background()
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-open-with-epoll",
	}

	mq, err := posixmq.OpenWithEpoll(ctx, cfg)

	if err != nil {
		t.Fatal(mq, err)
	}

	data := []byte("hello there!")
	// sending the bytes with priority ten.
	mq.Send(ctx, data, 10)
	mq.Close(ctx)
	mq.Unlink(ctx, cfg.Name)
}
