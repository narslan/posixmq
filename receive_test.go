package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestReceive(t *testing.T) {

	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-receive",
	}
	mq, err := posixmq.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}

	msg := "hello"
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
		t.Logf("response: %q", string(resp))
	}
	mq.Close()

	posixmq.Unlink(mq.Name)
}
