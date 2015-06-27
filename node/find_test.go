package node

import "testing"

func TestFindNode(t *testing.T) {
	// tests are [path -> ID]
	tests := map[string]string{
		"/":                   "/",
		"/README.md":          "/README.md",
		"/rEaDme.MD":          "/README.md",
		"//rEaDme.MD":         "/README.md",
		"///REadmE.Md":        "/README.md",
		"/pictuREs":           "/pictures",
		"/pictures/loGO.png":  "/pictures/logo.png",
		"/pictures//loGO.png": "/pictures/logo.png",
	}

	for path, ID := range tests {
		n, err := Mocked.FindNode(path)
		if err != nil {
			t.Fatalf("MockNodeTree.FindNode(%q) error: %s", path, err)
		}
		if want, got := ID, n.ID; want != got {
			t.Errorf("MockNodeTree.FindNode(%q).ID: want %s got %s", path, want, got)
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
