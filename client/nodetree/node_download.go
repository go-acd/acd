package nodetree

import (
	"fmt"
	"io"
	"net/http"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// Download downloads the node and returns the body as io.ReadCloser or an
// error. The caller is responsible for closing the reader.
func (n *Node) Download() (io.ReadCloser, error) {
	url := n.client.GetContentURL(fmt.Sprintf("nodes/%s/content", n.ID))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("error creating download request: %s", err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("error downloading the file: %s", err)
		return nil, constants.ErrDoingHTTPRequest
	}

	return res.Body, nil
}
