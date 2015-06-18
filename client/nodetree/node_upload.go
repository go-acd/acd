package nodetree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// CreateFolder creates the named folder under the node
func (n *Node) CreateFolder(name string) (*Node, error) {
	cn := &newNode{
		Name:    name,
		Kind:    "FOLDER",
		Parents: []string{n.ID},
	}
	jsonBytes, err := json.Marshal(cn)
	if err != nil {
		log.Errorf("error JSON encoding the metadata: %s", err)
		return nil, constants.ErrJSONEncoding
	}

	req, err := http.NewRequest("POST", n.client.GetMetadataURL("nodes"), bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("error creating the request to create the new folder: %s", err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("error creating the folder: %s", err)
		return nil, constants.ErrDoingHTTPRequest
	}
	if err := n.client.CheckResponse(res); err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var node Node
	if err := json.NewDecoder(res.Body).Decode(&node); err != nil {
		log.Errorf("error decoding the JSON body after creating the folder: %s", err)
		return nil, constants.ErrJSONDecodingResponseBody
	}
	n.AddChild(&node)

	return &node, nil
}

// Upload writes contents of r as name inside the current node.
func (n *Node) Upload(name string, r io.Reader) (*Node, error) {
	metadata := &newNode{
		Name:    name,
		Kind:    "FILE",
		Parents: []string{n.ID},
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Errorf("error JSON encoding the metadata: %s", err)
		return nil, constants.ErrJSONEncoding
	}

	postURL := n.client.GetContentURL("nodes?suppress=deduplication")
	node, err := n.upload(postURL, "POST", string(metadataJSON), name, r)
	if err != nil {
		return nil, err
	}

	n.AddChild(node)
	return node, nil
}

// Overwrite writes contents of r as name inside the current node.
func (n *Node) Overwrite(r io.Reader) error {
	putURL := n.client.GetContentURL(fmt.Sprintf("nodes/%s/content", n.ID))
	node, err := n.upload(putURL, "PUT", "", n.Name, r)
	if err != nil {
		return err
	}

	return n.update(node)
}

// TODO: Must go over this again
func (n *Node) upload(url, method, metadataJSON, name string, r io.Reader) (*Node, error) {
	bodyReader, bodyWriter := io.Pipe()
	writer := multipart.NewWriter(bodyWriter)
	errChan := make(chan error, 5)
	go func() {
		if metadataJSON != "" {
			if err := writer.WriteField("metadata", metadataJSON); err != nil {
				log.Errorf("%s: %s", constants.ErrWritingMetadata, err)
				errChan <- constants.ErrWritingMetadata
				return
			}
		}

		part, err := writer.CreateFormFile("content", name)
		if err != nil {
			log.Errorf("%s: %s", constants.ErrCreatingWriterFromFile, err)
			errChan <- err
			return
		}
		if _, err := io.Copy(part, r); err != nil {
			log.Errorf("%s: %s", constants.ErrWritingFileContents, err)
			errChan <- constants.ErrWritingFileContents
			return
		}

		errChan <- writer.Close()
		errChan <- bodyWriter.Close()
	}()

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		log.Errorf("error creating the upload request: %s", err)
		return nil, constants.ErrCreatingHTTPRequest
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("error uploading the file: %s", err)
		return nil, constants.ErrDoingHTTPRequest
	}
	if err := n.client.CheckResponse(res); err != nil {
		return nil, err
	}

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
	}

	defer res.Body.Close()
	var node Node
	if err := json.NewDecoder(res.Body).Decode(&node); err != nil {
		log.Errorf("error decoding the JSON body after uploading the file: %s", err)
		return nil, constants.ErrJSONDecodingResponseBody
	}

	return &node, nil
}
