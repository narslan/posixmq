package posixmq_test

import (
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

func TestReceive(t *testing.T) {

	m, err := posixmq.Open("test-receive")
	if err != nil {
		t.Fatal(m, err)
	}

	resp, err := m.Receive()
	if err != nil {
		t.Fatal(m, err)
	}

	t.Log("response:", string(resp))
}
