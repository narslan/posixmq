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
	m.Close()

}
