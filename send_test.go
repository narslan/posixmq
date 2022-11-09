package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestSend(t *testing.T) {

	qname := "test-send"
	ctx := context.Background()

	mq, err := posixmq.Open(ctx, qname, qsize, msize)
	if err != nil {
		t.Fatal(mq, err)
	}

	data := []byte("hello")

	err = mq.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(mq, err)
	}

	mq.Close(ctx)
	mq.Unlink(ctx, qname)
}

func BenchmarkSend(b *testing.B) {
	b.ReportAllocs()

	qname := "test-send"
	ctx := context.Background()

	mq, err := posixmq.Open(ctx, qname, qsize, msize)
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
	mq.Unlink(ctx, qname)
}
