package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestSendOpen(t *testing.T) {
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-send",
	}

	mq, err := posixmq.Open(cfg)
	if err != nil {
		t.Fatal(mq, err)
	}

	data := []byte("hello")

	err = mq.Send(data, 2)
	if err != nil {
		t.Fatal(mq, err)
	}

	mq.Close()
	posixmq.Unlink(mq.Name)
}

func TestSendCases(t *testing.T) {
	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-send1",
	}

	mq, err := posixmq.Open(cfg)
	if err != nil {
		t.Fatal(mq, err)
	}

	for i := 0; i < int(cfg.QueueSize); i++ {
		data := []byte("hello")
		err = mq.Send(data, 2)

		if err != nil {

			err := mq.Close()
			if err != nil {
				t.Fatal(err)
			}
			posixmq.Unlink(mq.Name)
			if err != nil {
				t.Fatal(err)
			}

		}
	}

	mq.Close()
	posixmq.Unlink(mq.Name)

}

func BenchmarkSend(b *testing.B) {
	b.ReportAllocs()

	cfg := &posixmq.Config{
		QueueSize:   10,
		MessageSize: 4096,
		Name:        "test-unlink",
	}
	mq, err := posixmq.Open(cfg)
	if err != nil {
		b.Fatal(mq, err)
	}

	for n := 0; n < b.N; n++ {

		data := []byte("hello")
		err = mq.Send(data, 2)
		if err != nil {
			b.Fatal(mq, err)
		}
	}
	mq.Close()
	posixmq.Unlink(mq.Name)

}
