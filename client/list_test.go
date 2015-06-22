package client

import (
	"reflect"
	"testing"

	"gopkg.in/acd.v0/client/nodetree"
)

func TestList(t *testing.T) {
	c := &Client{
		NodeTree: nodetree.Mocked,
	}

	tests := map[string][]string{
		"/":         []string{"README.md", "pictures"},
		"/pictures": []string{"logo.png"},
	}

	for path, want := range tests {
		var names []string
		nodes, err := c.List(path)
		if err != nil {
			t.Errorf("c.List(%q) error: %s", path, err)
		}
		for _, node := range nodes {
			names = append(names, node.Name)
		}
		if got := names; !reflect.DeepEqual(want, got) {
			t.Errorf("c.List(%q): want %v got %v", "/", want, got)
		}
	}
}
