package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestSendOpen(t *testing.T) {
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-send",
	}

	ctx := context.Background()
	mq, err := posixmq.Open(ctx, cfg)
	if err != nil {
		t.Fatal(mq, err)
	}

	data := []byte("hello")

	err = mq.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(mq, err)
	}

	mq.Close(ctx)
	mq.Unlink(ctx)
}

func TestSendCases(t *testing.T) {
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-send1",
	}

	ctx := context.Background()
	mq, err := posixmq.Open(ctx, cfg)
	if err != nil {
		t.Fatal(mq, err)
	}

	for i := 0; i < int(cfg.QueueSize); i++ {
		data := []byte("hello")
		err = mq.Send(ctx, data, 2)

		if err != nil {

			err := mq.Close(ctx)
			if err != nil {
				t.Fatal(err)
			}
			err = mq.Unlink(ctx)
			if err != nil {
				t.Fatal(err)
			}

		}
	}

	mq.Close(ctx)
	mq.Unlink(ctx)
}

func BenchmarkSend(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()

	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-unlink",
	}
	mq, err := posixmq.Open(ctx, cfg)
	if err != nil {
		b.Fatal(mq, err)
	}

	for n := 0; n < b.N; n++ {

		data := []byte("hello")
		err = mq.Send(ctx, data, 2)
		if err != nil {
			b.Fatal(mq, err)
		}
	}
	mq.Close(ctx)
	mq.Unlink(ctx)
}
