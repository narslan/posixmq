package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestReceive(t *testing.T) {

	msg := "hello"
	qname := "test-receive"
	mq, err := posixmq.Open(qname)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte(msg)
	err = mq.Send(data, 0)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := mq.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if string(resp) != msg {
		t.Log("response:", string(resp))
	}
	err = mq.Unlink(qname)
	if err != nil {
		t.Fatal(mq, err)
	}
}
