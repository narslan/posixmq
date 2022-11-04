package posixmq_test

import (
	"testing"

	"github.com/narslan/posixmq"
)

func TestSend(t *testing.T) {

	qname := "test-send"
	m, err := posixmq.Open(qname)
	if err != nil {
		t.Fatal(m, err)
	}

	data := []byte("hello")
	err = m.Send(data, 2)
	if err != nil {
		t.Fatal(m, err)
	}
}
