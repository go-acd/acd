package node

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

type (
	changes struct {
		Checkpoint    string `json:"checkpoint,omitempty"`
		Chunksize     int    `json:"chunkSize,omitempty"`
		MaxNodes      int    `json:"maxNodes,omitempty"`
		IncludePurged string `json:"includePurged,omitempty"`
	}

	changesResponse struct {
		Checkpoint string  `json:"checkpoint,omitempty"`
		Nodes      []*Node `json:"nodes,omitempty"`
		Reset      bool    `json:"reset,omitempty"`
		End        bool    `json:"end,omitempty"`
	}
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

	// return format should be:
	// {"checkpoint": str, "reset": bool, "nodes": []}
	// {"checkpoint": str, "reset": false, "nodes": []}
	// {"end": true}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrReadingResponseBody, err)
		return constants.ErrReadingResponseBody
	}
	for _, lineBytes := range bytes.Split(bodyBytes, []byte("\n")) {
		var cr changesResponse
		if err := json.Unmarshal(lineBytes, &cr); err != nil {
			log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
			return constants.ErrJSONDecodingResponseBody
		}
		if cr.Checkpoint != "" {
			log.Debugf("changes returned Checkpoint: %s", cr.Checkpoint)
			nt.Checkpoint = cr.Checkpoint
		}
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
		log.Debugf("node %s ID %s has changed.", node.Name, node.ID)
		if _, found := nt.nodeMap[node.ID]; !found {
			// remove the parents from the node we are inserting so the next section
			// will detect the added parents and add those.
			newNode := &Node{}
			(*newNode) = *node
			newNode.Parents = []string{}
			nt.nodeMap[node.ID] = newNode
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
				parent.RemoveChild(nt.nodeMap[node.ID])
			}

			// remove the node itself from the nodemap
			delete(nt.nodeMap, node.ID)

			continue
		}

		// add/remove parents
		sort.Strings(nt.nodeMap[node.ID].Parents)
		sort.Strings(newNode.Parents)
		if parentIDs := diffSliceStr(nt.nodeMap[node.ID].Parents, newNode.Parents); len(parentIDs) > 0 {
			for _, parentID := range parentIDs {
				log.Debugf("ParentID %s has been removed from %s ID %s", parentID, node.Name, node.ID)
				parent, err := nt.FindByID(parentID)
				if err != nil {
					continue
				}
				parent.RemoveChild(nt.nodeMap[node.ID])
			}
		}
		if parentIDs := diffSliceStr(newNode.Parents, nt.nodeMap[node.ID].Parents); len(parentIDs) > 0 {
			for _, parentID := range parentIDs {
				log.Debugf("ParentID %s has been added to %s ID %s", parentID, node.Name, node.ID)
				parent, err := nt.FindByID(parentID)
				if err != nil {
					continue
				}
				parent.AddChild(nt.nodeMap[node.ID])
			}
		}

		// finally update the node itself
		(*node) = *newNode
	}

	return nil
}
