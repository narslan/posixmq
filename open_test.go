package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestOpen(t *testing.T) {

	qname := "test-open"
	ctx := context.Background()
	m, err := posixmq.Open(ctx, qname)
	if err != nil {
		t.Fatal(m, err)
	}
	m.Close(ctx)

}

func TestUnlink(t *testing.T) {

	msg := "hello"
	qname := "test-unlink"
	ctx := context.Background()
	mq, err := posixmq.Open(ctx, qname)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte(msg)
	err = mq.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = mq.Unlink(ctx, qname)
	if err != nil {
		t.Fatal(mq)
	}
}
