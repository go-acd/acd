package client

import (
	"gopkg.in/acd.v0/client/node"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// List returns a list of os.Fileinfo representing the files under the path
func (c *Client) List(path string) (node.Nodes, error) {
	rootNode, err := c.GetNodeTree().FindNode(path)
	if err != nil {
		return nil, err
	}
	if !rootNode.IsDir() {
		log.Errorf("%s: %s", constants.ErrPathIsNotFolder, path)
		return nil, constants.ErrPathIsNotFolder
	}

	return rootNode.Nodes, nil
}
