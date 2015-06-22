package node

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
		log.Errorf("%s: %s", constants.ErrJSONEncoding, err)
		return nil, constants.ErrJSONEncoding
	}

	req, err := http.NewRequest("POST", n.client.GetMetadataURL("nodes"), bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return nil, constants.ErrCreatingHTTPRequest
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := n.client.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return nil, constants.ErrDoingHTTPRequest
	}
	if err := n.client.CheckResponse(res); err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var node Node
	if err := json.NewDecoder(res.Body).Decode(&node); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
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
		log.Errorf("%s: %s", constants.ErrJSONEncoding, err)
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

func (n *Node) upload(url, method, metadataJSON, name string, r io.Reader) (*Node, error) {
	bodyReader, bodyWriter := io.Pipe()
	errChan := make(chan error)
	bodyChan := make(chan io.ReadCloser)
	contentTypeChan := make(chan string)

	go n.bodyWriter(metadataJSON, name, r, bodyWriter, errChan, contentTypeChan)
	go func() {
		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
			select {
			case errChan <- constants.ErrCreatingHTTPRequest:
			default:
			}
			return
		}
		req.Header.Add("Content-Type", <-contentTypeChan)
		res, err := n.client.Do(req) // this should block until the upload is finished.
		if err != nil {
			log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
			select {
			case errChan <- constants.ErrDoingHTTPRequest:
			default:
			}
			return
		}
		if err := n.client.CheckResponse(res); err != nil {
			select {
			case errChan <- err:
			default:
			}
			return
		}

		select {
		case bodyChan <- res.Body:
		default:
		}
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		case body := <-bodyChan:
			defer body.Close()
			var node Node
			if err := json.NewDecoder(body).Decode(&node); err != nil {
				log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
				return nil, constants.ErrJSONDecodingResponseBody
			}

			return &node, nil
		}
	}
}

func (n *Node) bodyWriter(metadataJSON, name string, r io.Reader, bodyWriter io.WriteCloser, errChan chan error, contentTypeChan chan string) {
	writer := multipart.NewWriter(bodyWriter)
	contentTypeChan <- writer.FormDataContentType()
	if metadataJSON != "" {
		if err := writer.WriteField("metadata", metadataJSON); err != nil {
			log.Errorf("%s: %s", constants.ErrWritingMetadata, err)
			select {
			case errChan <- constants.ErrWritingMetadata:
			default:
			}
			return
		}
	}

	part, err := writer.CreateFormFile("content", name)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingWriterFromFile, err)
		select {
		case errChan <- err:
		default:
		}
		return
	}
	if _, err := io.Copy(part, r); err != nil {
		log.Errorf("%s: %s", constants.ErrWritingFileContents, err)
		select {
		case errChan <- constants.ErrWritingFileContents:
		default:
		}
		return
	}

	select {
	case errChan <- writer.Close():
	default:
	}
	select {
	case errChan <- bodyWriter.Close():
	default:
	}
}
