package node

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
	if n.IsDir() {
		log.Errorf("%s: cannot download a folder", constants.ErrPathIsFolder)
		return nil, constants.ErrPathIsFolder
	}
	url := n.client.GetContentURL(fmt.Sprintf("nodes/%s/content", n.ID))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return nil, constants.ErrCreatingHTTPRequest
	}
	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return nil, constants.ErrDoingHTTPRequest
	}
	if err := n.client.CheckResponse(res); err != nil {
		return nil, err
	}

	return res.Body, nil
}
