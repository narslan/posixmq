package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestReceive(t *testing.T) {

	m, err := posixmq.Open("test-receive")
	if err != nil {
		t.Fatal(m, err)
	}

	data := []byte("hello")
	err = m.Send(data, 0)
	if err != nil {
		t.Fatal(m, err)
	}

	resp, err := m.Receive()
	if err != nil {
		t.Fatal(m, err)
	}

	t.Log("response:", string(resp))
}
