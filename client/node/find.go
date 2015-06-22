package node

import (
	"os"
	"strings"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// FindNode finds a node for a particular path.
// TODO(kalbasit): This does not perform well, this should be cached in a map
// path->node and calculated on load (fresh, cache, refresh).
func (nt *Tree) FindNode(path string) (*Node, error) {
	// chop off the first PathSeparator.
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		path = path[1:]
	}

	// initialize our search from the root node
	node := nt.Node

	// are we looking for the root node?
	if path == "" {
		return node, nil
	}

	// iterate over the path parts until we find the path (or not).
	parts := strings.Split(path, "/")
	for _, part := range parts {
		var found bool
		for _, n := range node.Nodes {
			// does node.name matches our query?
			if n.Name == part {
				node = n
				found = true
				break
			}
		}

		if !found {
			log.Errorf("%s: %s", constants.ErrNodeNotFound, path)
			return nil, constants.ErrNodeNotFound
		}
	}

	return node, nil
}

// FindByID returns the node identified by the ID.
func (nt *Tree) FindByID(id string) (*Node, error) {
	n, found := nt.nodeMap[id]
	if !found {
		log.Errorf("%s: ID %q", constants.ErrNodeNotFound, id)
		return nil, constants.ErrNodeNotFound
	}

	return n, nil
}
