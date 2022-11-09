package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

var qsize int64 = 8    // number of message that queue accepts
var msize int64 = 4096 // limit of size of each message

func TestOpen(t *testing.T) {

	qname := "test-open"
	ctx := context.Background()
	mq, err := posixmq.Open(ctx, qname, qsize, msize)

	if err != nil {
		t.Fatal(mq, err)
	}
	mq.Close(ctx)

}

func TestUnlink(t *testing.T) {

	msg := "hello"
	qname := "test-unlink"

	ctx := context.Background()
	mq, err := posixmq.Open(ctx, qname, qsize, msize)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte(msg)
	err = mq.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(err)
	}

	mq.Close(ctx)
	err = mq.Unlink(ctx, qname)
	if err != nil {
		t.Fatal(mq)
	}
}
