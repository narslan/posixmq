package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestReceive(t *testing.T) {

	msg := "hello"
	qname := "test-receive"
	ctx := context.Background()

	mq, err := posixmq.Open(ctx, qname, qsize, msize)

	if err != nil {
		t.Fatal(err)
	}

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
	err = mq.Unlink(ctx, qname)
	if err != nil {
		t.Fatal(mq, err)
	}
}
