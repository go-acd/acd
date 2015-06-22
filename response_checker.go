package acd

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// CheckResponse validates the response from the Amazon Cloud Drive API. It
// does that by looking at the response's status code and it returns an error
// for any code lower than 200 or greater than 300
func (c *Client) CheckResponse(res *http.Response) error {
	if 200 <= res.StatusCode && res.StatusCode <= 299 {
		return nil
	}
	errBody := "no response body"
	defer res.Body.Close()
	if data, err := ioutil.ReadAll(res.Body); err == nil {
		errBody = string(data)
	}
	var err error
	switch res.StatusCode {
	case http.StatusBadRequest:
		err = constants.ErrResponseBadInput
	case http.StatusUnauthorized:
		err = constants.ErrResponseInvalidToken
	case http.StatusForbidden:
		err = constants.ErrResponseForbidden
	case http.StatusConflict:
		err = constants.ErrResponseDuplicateExists
	case http.StatusInternalServerError:
		err = constants.ErrResponseInternalServerError
	case http.StatusServiceUnavailable:
		err = constants.ErrResponseUnavailable
	default:
		err = constants.ErrResponseUnknown
	}

	log.Errorf("{code: %s} %s: %s", res.Status, err, errBody)
	return err
}
