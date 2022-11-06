package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestSend(t *testing.T) {

	qname := "test-send"
	ctx := context.Background()
	m, err := posixmq.Open(ctx, qname)
	if err != nil {
		t.Fatal(m, err)
	}

	data := []byte("hello")
	err = m.Send(ctx, data, 2)
	if err != nil {
		t.Fatal(m, err)
	}
}
