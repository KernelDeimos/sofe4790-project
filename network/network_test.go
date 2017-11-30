package network

import (
	"testing"
)

func TestNodeSendGeneratesCorrectURI(t *testing.T) {
	node := &PeerNode{
		Host: "127.0.0.1",
		Port: 3000,
	}

	uri, err := node.generateURI("_event", struct {
		Attr1 string
		Attr2 string
	}{"test", "values"})

	if err != nil {
		t.Error(err)
	}

	if uri != "http://127.0.0.1:3000/_event/%7B%22Attr1%22:%22test%22%2C%22Attr2%22:%22values%22%7D" {
		t.Fail()
	}
}
