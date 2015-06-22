package nodetree

import "testing"

func TestFindNode(t *testing.T) {
	tests := []string{
		"/",
		"/README.md",
		"/pictures",
		"/pictures/logo.png",
	}

	for _, test := range tests {
		n, err := Mocked.FindNode(test)
		if err != nil {
			t.Errorf("MockNodeTree.FindNode(%q) error: %s", test, err)
		}
		if want, got := test, n.ID; want != got {
			t.Errorf("MockNodeTree.FindNode(%q).ID: want %s got %s", test, want, got)
		}
	}
}

func TestFindById(t *testing.T) {
	tests := []string{
		"/",
		"/README.md",
		"/pictures",
		"/pictures/logo.png",
	}

	for _, test := range tests {
		n, err := Mocked.FindByID(test)
		if err != nil {
			t.Errorf("MockNodeTree.FindByID(%q) error: %s", test, err)
		}
		if want, got := test, n.ID; want != got {
			t.Errorf("MockNodeTree.FindByID(%q).ID: want %s got %s", test, want, got)
		}
	}
}
