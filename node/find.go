package node

import (
	"regexp"
	"strings"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// FindNode finds a node for a particular path.
// TODO(kalbasit): This does not perform well, this should be cached in a map
// path->node and calculated on load (fresh, cache, refresh).
func (nt *Tree) FindNode(path string) (*Node, error) {
	// replace multiple n*/ with /
	re := regexp.MustCompile("/[/]*")
	path = string(re.ReplaceAll([]byte(path), []byte("/")))
	// chop off the first /.
	path = strings.TrimPrefix(path, "/")
	// did we ask for the root node?
	if path == "" {
		return nt.Node, nil
	}

	// initialize our search from the root node
	node := nt.Node

	// iterate over the path parts until we find the path (or not).
	parts := strings.Split(path, "/")
	for _, part := range parts {
		var found bool
		for _, n := range node.Nodes {
			// does node.name matches our query?
			if strings.ToLower(n.Name) == strings.ToLower(part) {
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
