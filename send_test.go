package posixmq_test

import (
	"context"
	"testing"

	"github.com/narslan/posixmq"
)

func TestOpen(t *testing.T) {

	m, err := posixmq.Open("test")
	if err != nil {
		t.Fatal(m, err)
	}

	m, err = posixmq.Open("/")
	if err == nil {
		t.Fatal(m, err)
	}

}

func TestSend(t *testing.T) {

	m, err := posixmq.Open("test")
	if err != nil {
		t.Fatal(m, err)
	}

	ctx := context.Background()
	data := []byte("hello")
	err = m.Send(ctx, data)
	if err != nil {
		t.Fatal(m, err)
	}
}
