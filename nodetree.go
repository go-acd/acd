package acd

import "gopkg.in/acd.v0/node"

// FetchNodeTree fetches and caches the NodeTree.
func (c *Client) FetchNodeTree() error {
	nt, err := node.NewTree(c, c.cacheFile)
	if err != nil {
		return err
	}

	c.NodeTree = nt
	return nil
}

// GetNodeTree returns the NodeTree.
func (c *Client) GetNodeTree() *node.Tree {
	return c.NodeTree
}
