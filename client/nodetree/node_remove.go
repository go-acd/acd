package nodetree

import (
	"fmt"
	"net/http"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// Remove deletes a node from the server.
// This function does not update the NodeTree, the caller should do so!
func (n *Node) Remove() error {
	putURL := n.client.GetMetadataURL(fmt.Sprintf("/trash/%s", n.ID))

	req, err := http.NewRequest("PUT", putURL, nil)
	if err != nil {
		log.Errorf("error creating the trash request: %s", err)
		return constants.ErrCreatingHTTPRequest
	}

	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("error deleting the file: %s", err)
		return constants.ErrDoingHTTPRequest
	}
	if err := n.client.CheckResponse(res); err != nil {
		return err
	}

	return nil
}
