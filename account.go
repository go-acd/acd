package acd

import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

type (
	// AccountInfo represents information about an Amazon Cloud Drive account.
	AccountInfo struct {
		TermsOfUse string `json:"termsOfUse"`
		Status     string `json:"status"`
	}

	// AccountQuota represents information about the account quotas.
	AccountQuota struct {
		Quota          uint64    `json:"quota"`
		LastCalculated time.Time `json:"lastCalculated"`
		Available      uint64    `json:"available"`
	}

	// AccountUsage represents information about the account usage.
	AccountUsage struct {
		LastCalculated time.Time `json:"lastCalculated"`

		Doc struct {
			Billable struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"billable"`
			Total struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"total"`
		} `json:"doc"`

		Other struct {
			Billable struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"billable"`
			Total struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"total"`
		} `json:"other"`

		Photo struct {
			Billable struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"billable"`
			Total struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"total"`
		} `json:"photo"`

		Video struct {
			Billable struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"billable"`
			Total struct {
				Bytes uint64 `json:"bytes"`
				Count uint32 `json:"count"`
			} `json:"total"`
		} `json:"video"`
	}
)

// GetAccountInfo returns AccountInfo about the current account.
func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	var ai AccountInfo
	req, err := http.NewRequest("GET", c.metadataURL+"/account/info", nil)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	res, err := c.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return nil, constants.ErrDoingHTTPRequest
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&ai); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
		return nil, constants.ErrJSONDecodingResponseBody
	}

	return &ai, nil
}

// GetAccountQuota returns AccountQuota about the current account.
func (c *Client) GetAccountQuota() (*AccountQuota, error) {
	var aq AccountQuota
	req, err := http.NewRequest("GET", c.metadataURL+"/account/quota", nil)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	res, err := c.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return nil, constants.ErrDoingHTTPRequest
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&aq); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
		return nil, constants.ErrJSONDecodingResponseBody
	}

	return &aq, nil
}

// GetAccountUsage returns AccountUsage about the current account.
func (c *Client) GetAccountUsage() (*AccountUsage, error) {
	var au AccountUsage
	req, err := http.NewRequest("GET", c.metadataURL+"/account/usage", nil)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	res, err := c.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return nil, constants.ErrDoingHTTPRequest
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&au); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
		return nil, constants.ErrJSONDecodingResponseBody
	}

	return &au, nil
}
