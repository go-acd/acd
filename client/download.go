package client

import "io"

// Download returns an io.ReadCloser for the file specified by the path. The
// caller is responsible for closing the body.
func (c *Client) Download(path string) (io.ReadCloser, error) {
	node, err := c.NodeTree.FindNode(path)
	if err != nil {
		return nil, err
	}

	return node.Download()
}
