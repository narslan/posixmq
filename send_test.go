package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestSend(t *testing.T) {

	m, err := posixmq.Open("test-send")
	if err != nil {
		t.Fatal(m, err)
	}

	data := []byte("hello")
	err = m.Send(data, 2)
	if err != nil {
		t.Fatal(m, err)
	}
}

func TestUnlink(t *testing.T) {

	m, err := posixmq.Open("test-unlink")
	if err != nil {
		t.Fatal(m, err)
	}

	data := []byte("hello")
	err = m.Send(data, 2)
	if err != nil {
		t.Fatal(m, err)
	}
	err = m.Unlink("test-unlink")
	if err != nil {
		t.Fatal(m, err)
	}
}
