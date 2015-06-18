package client

import "gopkg.in/acd.v0/client/nodetree"

// FetchNodeTree fetches and caches the NodeTree
func (c *Client) FetchNodeTree() error {
	nt, err := nodetree.New(c, c.cacheFile)
	if err != nil {
		return err
	}

	c.NodeTree = nt
	return nil
}

// GetNodeTree returns the NodeTree.
func (c *Client) GetNodeTree() *nodetree.NodeTree {
	return c.NodeTree
}
