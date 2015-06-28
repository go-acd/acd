// Package token represents an oauth2.TokenSource which has the ability to
// refresh the access token through the oauth server.
package token // import "gopkg.in/acd.v0/token"

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

const refreshURL = "https://go-acd.appspot.com/refresh"

// Source provides a Source with support for refreshing from the acd server.
type Source struct {
	path  string
	token *oauth2.Token
}

// New returns a new Source implementing oauth2.TokenSource. The path must
// exist on the filesystem and must be of permissions 0600.
func New(path string) (*Source, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Errorf("%s: %s", constants.ErrFileNotFound, path)
		return nil, constants.ErrFileNotFound
	}

	ts := &Source{
		path:  path,
		token: new(oauth2.Token),
	}
	ts.readToken()

	return ts, nil
}

// Token returns an oauth2.Token. If the cached token (in (*Source).path) has
// expired, it will fetch the token from the server and cache it before
// returning it.
func (ts *Source) Token() (*oauth2.Token, error) {
	if !ts.token.Valid() {
		log.Debug("token is not valid, it has probably expired")
		if err := ts.refreshToken(); err != nil {
			return nil, err
		}

		if err := ts.saveToken(); err != nil {
			return nil, err
		}
	}

	return ts.token, nil
}

func (ts *Source) readToken() error {
	log.Debugf("reading the token from %s", ts.path)
	f, err := os.Open(ts.path)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrOpenFile, ts.path)
		return constants.ErrOpenFile
	}
	if err := json.NewDecoder(f).Decode(ts.token); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecoding, err)
		return constants.ErrJSONDecoding
	}

	log.Debug("token loaded successfully")
	return nil
}

func (ts *Source) saveToken() error {
	log.Debugf("saving the token to %s", ts.path)
	f, err := os.Create(ts.path)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreateFile, ts.path)
		return constants.ErrCreateFile
	}
	if err := json.NewEncoder(f).Encode(ts.token); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONEncoding, err)
		return constants.ErrJSONEncoding
	}

	log.Debug("token saved successfully")
	return nil
}

func (ts *Source) refreshToken() error {
	log.Debugf("refreshing the token from %q", refreshURL)

	data, err := json.Marshal(ts.token)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrJSONEncoding, err)
		return constants.ErrJSONEncoding
	}
	req, err := http.NewRequest("POST", refreshURL, bytes.NewBuffer(data))
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return constants.ErrCreatingHTTPRequest
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return constants.ErrDoingHTTPRequest
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(ts.token); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
		return constants.ErrJSONDecodingResponseBody
	}
	log.Debug("token was refreshed successfully")

	return nil
}
