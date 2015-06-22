package acd

import (
	"encoding/json"
	"net/http"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// GetMetadataURL returns the metadata url.
func (c *Client) GetMetadataURL(path string) string {
	return c.metadataURL + path
}

// GetContentURL returns the content url.
func (c *Client) GetContentURL(path string) string {
	return c.contentURL + path
}

func setEndpoints(c *Client) error {
	req, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return constants.ErrCreatingHTTPRequest
	}

	var er endpointResponse
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return constants.ErrDoingHTTPRequest
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
		return constants.ErrJSONDecodingResponseBody
	}

	c.contentURL = er.ContentURL
	c.metadataURL = er.MetadataURL
	return nil
}
