package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestOpen(t *testing.T) {

	qname := "test-open"
	m, err := posixmq.Open(qname)
	if err != nil {
		t.Fatal(m, err)
	}
	m.Close()

}

func TestUnlink(t *testing.T) {

	msg := "hello"
	qname := "test-unlink"

	mq, err := posixmq.Open(qname)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte(msg)
	err = mq.Send(data, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = mq.Unlink(qname)
	if err != nil {
		t.Fatal(mq)
	}
}
