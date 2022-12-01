package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

// var qsize int64 = 8    // number of message that queue accepts
// var msize int64 = 4096 // limit of size of each message
func TestOpen(t *testing.T) {
	ctx := context.Background()
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-open",
	}

	mq, err := posixmq.Open(ctx, cfg)
	if err != nil {
		t.Fatal(mq, err)
	}
	t.Log("-------")
	t.Log(mq)
	t.Log(cfg)

	mq.Close(ctx)
	mq.Unlink(ctx, cfg.Name)
}

func TestUnlink(t *testing.T) {
	ctx := context.Background()
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 1024,
		Name:        "test-unlink",
	}
	mq, err := posixmq.Open(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("test byte")
	// with priority 2
	err = mq.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(err)
	}

	mq.Close(ctx)
	err = mq.Unlink(ctx, cfg.Name)
	if err != nil {
		t.Fatal(mq)
	}
}
