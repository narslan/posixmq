package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

// var qsize int64 = 8    // number of message that queue accepts
// var msize int64 = 4096 // limit of size of each message
func TestOpen(t *testing.T) {

	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-open",
	}

	mq, err := posixmq.Open(cfg)
	if err != nil {
		t.Fatal(mq, err)
	}

	mq.Close()
	mq.Unlink()
}

func TestUnlink(t *testing.T) {

	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 1024,
		Name:        "test-unlink",
	}
	mq, err := posixmq.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("test byte")
	// with priority 2
	err = mq.Send(data, 2)
	if err != nil {
		t.Fatal(err)
	}

	mq.Close()
	err = mq.Unlink()
	if err != nil {
		t.Fatal(mq)
	}
}
