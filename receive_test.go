package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestReceive(t *testing.T) {

	ctx := context.Background()

	cfg := &posixmq.Config{
		QueueSize:   100,
		MessageSize: 100,
		Name:        "test-receive",
	}
	mq, err := posixmq.Open(ctx, cfg)

	if err != nil {
		t.Fatal(err)
	}

	msg := "hello"
	data := []byte(msg)
	err = mq.Send(ctx, data, 0)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := mq.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if string(resp) != msg {
		t.Log("response:", string(resp))
	}
}
