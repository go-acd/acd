package node

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"sort"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// Sync syncs the tree with the server.
func (nt *Tree) Sync() error {
	postURL := nt.client.GetMetadataURL("changes")
	c := &changes{
		Checkpoint: nt.Checkpoint,
	}
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrJSONEncoding, err)
		return constants.ErrJSONEncoding
	}

	// return format should be:
	// {"checkpoint": str, "reset": bool, "nodes": []}
	// {"checkpoint": str, "reset": false, "nodes": []}
	// {"end": true}
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
		return constants.ErrCreatingHTTPRequest
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := nt.client.Do(req)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
		return constants.ErrDoingHTTPRequest
	}
	if err := nt.client.CheckResponse(res); err != nil {
		return err
	}

	defer res.Body.Close()
	s := bufio.NewScanner(res.Body)
	for s.Scan() {
		var cr changesResponse
		if err := json.Unmarshal(s.Bytes(), &cr); err != nil {
			log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
			return constants.ErrJSONDecodingResponseBody
		}
		nt.Checkpoint = cr.Checkpoint
		if cr.Reset {
			log.Debug("reset is required")
			return constants.ErrMustFetchFresh
		}
		if cr.End {
			break
		}
		if err := nt.updateNodes(cr.Nodes); err != nil {
			return err
		}
	}

	return nil
}

func (nt *Tree) updateNodes(nodes []*Node) error {
	// first make sure our nodeMap is up to date
	for _, node := range nodes {
		if _, found := nt.nodeMap[node.ID]; !found {
			nt.nodeMap[node.ID] = node
		}
	}

	// now let's update the nodes
	for _, node := range nodes {
		// make a copy of n
		newNode := &Node{}
		(*newNode) = *nt.nodeMap[node.ID]
		if err := newNode.update(node); err != nil {
			return err
		}

		// has this node been deleted?
		if !newNode.Available() {
			log.Debugf("node ID %s name %s has been deleted", newNode.ID, newNode.Name)
			for _, parentID := range append(newNode.Parents, nt.nodeMap[node.ID].Parents...) {
				parent, err := nt.FindByID(parentID)
				if err != nil {
					continue
				}
				parent.RemoveChild(node)
			}

			// remove the node itself from the nodemap
			delete(nt.nodeMap, node.ID)

			continue
		}

		// TODO: Handle change in parent IDs gracefully!
		sort.Strings(node.Parents)
		sort.Strings(newNode.Parents)
		if !reflect.DeepEqual(node.Parents, newNode.Parents) {
			log.Debugf("The parents of the node %s have changed. %v => %v", node.ID, node.Parents, newNode.Parents)
			return constants.ErrMustFetchFresh
		}

		// finally update the node itself
		(*node) = *newNode
	}

	return nil
}
