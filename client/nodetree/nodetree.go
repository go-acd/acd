package nodetree

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

type (
	// NodeTree points to the root node in the node tree
	NodeTree struct {
		*Node

		// Internal
		LastUpdated time.Time
		Checkpoint  string

		client    client
		cacheFile string
		nodeMap   map[string]*Node
	}

	nodeList struct {
		ETagResponse string  `json:"eTagResponse"`
		Count        uint64  `json:"count,omitempty"`
		NextToken    string  `json:"nextToken,omitempty"`
		Nodes        []*Node `json:"data,omitempty"`
	}

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

// RemoveNode removes this node from the server and from the NodeTree.
func (nt *NodeTree) RemoveNode(n *Node) error {
	if err := n.Remove(); err != nil {
		return err
	}

	for _, parentID := range n.Parents {
		parent, err := nt.FindByID(parentID)
		if err != nil {
			if err != constants.ErrNodeNotFound {
				log.Errorf("error trying to get the parent with ID %s: %s", parentID, err)
			}

			continue
		}
		parent.RemoveChild(n)
	}

	return nil
}

func (nt *NodeTree) loadOrRefresh() error {
	var err error
	if err = nt.loadCache(); err != nil {
		log.Debug(err)
		// first fetch the checkpoint by making a refresh request
		if err = nt.refresh(); err != nil && err != constants.ErrMustRefresh {
			return err
		}
		// now fetch a fresh copy
		if err = nt.fetchFresh(); err != nil {
			return err
		}
	}
	nt.buildNodeMap(nt.Node)

	if err = nt.refresh(); err != nil {
		switch err {
		case constants.ErrMustRefresh:
			if err := nt.removeCache(); err != nil && os.IsNotExist(err) {
				return err
			}
			return nt.loadOrRefresh()
		default:
			err = fmt.Errorf("error refreshing the NodeTree: %s", err)
			return err
		}
	}

	return nil
}

func (nt *NodeTree) removeCache() error {
	return os.Remove(nt.cacheFile)
}

// New returns the root node (the head of the tree).
func New(c client, cacheFile string) (*NodeTree, error) {
	nt := &NodeTree{
		cacheFile: cacheFile,
		client:    c,
	}
	if err := nt.loadOrRefresh(); err != nil {
		return nil, err
	}
	if err := nt.saveCache(); err != nil {
		return nil, err
	}

	return nt, nil
}

// Close finalizes the NodeTree
func (nt *NodeTree) Close() error {
	return nt.saveCache()
}

// FindByID returns the node identified by the ID.
func (nt *NodeTree) FindByID(id string) (*Node, error) {
	n, found := nt.nodeMap[id]
	if !found {
		log.Errorf("%s: ID %q", constants.ErrNodeNotFound, id)
		return nil, constants.ErrNodeNotFound
	}

	return n, nil
}

func (nt *NodeTree) fetchFresh() error {
	// grab the list of all of the nodes from the server.
	var nextToken string
	var nodes []*Node
	for {
		nl := nodeList{
			Nodes: make([]*Node, 0, 200),
		}
		urlStr := nt.client.GetMetadataURL("nodes")
		u, err := url.Parse(urlStr)
		if err != nil {
			log.Errorf("%s: %s", constants.ErrParsingURL, urlStr)
			return constants.ErrParsingURL
		}

		v := url.Values{}
		v.Set("limit", "200")
		if nextToken != "" {
			v.Set("startToken", nextToken)
		}
		u.RawQuery = v.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			log.Errorf("error creating the NodeTree request: %s", err)
			return constants.ErrCreatingHTTPRequest
		}

		req.Header.Set("Content-Type", "application/json")
		res, err := nt.client.Do(req)
		if err != nil {
			log.Errorf("error fetching the list of nodes: %s", err)
			return constants.ErrDoingHTTPRequest
		}
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&nl); err != nil {
			log.Errorf("error decoding the request body: %s", err)
			return constants.ErrJSONDecodingResponseBody
		}

		nextToken = nl.NextToken
		nodes = append(nodes, nl.Nodes...)

		if nextToken == "" {
			break
		}
	}

	nodeMap := make(map[string]*Node, len(nodes))
	for _, node := range nodes {
		if !node.Available() {
			continue
		}
		node.client = nt.client
		nodeMap[node.ID] = node
	}

	for _, node := range nodeMap {
		if node.Name == "" && node.IsFolder() && len(node.Parents) == 0 {
			nt.Node = node
			node.Root = true
		}

		for _, parentID := range node.Parents {
			if pn, found := nodeMap[parentID]; found {
				pn.Nodes = append(pn.Nodes, node)
			}
		}
	}

	nt.nodeMap = nodeMap
	return nil
}

func (nt *NodeTree) buildNodeMap(current *Node) {
	if nt.Node == current {
		nt.nodeMap = make(map[string]*Node)
	}
	nt.nodeMap[current.ID] = current
	for _, node := range current.Nodes {
		nt.buildNodeMap(node)
	}
}

// FindNode finds a node for a particular path.
func (nt *NodeTree) FindNode(path string) (*Node, error) {
	// chop off the first PathSeparator.
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		path = path[1:]
	}

	// initialize our search from the root node
	node := nt.Node

	// are we looking for the root node?
	if path == "" {
		return node, nil
	}

	// iterate over the path parts until we find the path (or not).
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		var found bool
		for _, n := range node.Nodes {
			// does node.name matches our query?
			if n.Name == part {
				node = n
				found = true
				break
			}
		}

		if !found {
			log.Errorf("%s: %s", constants.ErrNodeNotFound, path)
			return nil, constants.ErrNodeNotFound
		}
	}

	return node, nil
}

// CreatePath creates a folder under the given path.
func (nt *NodeTree) CreatePath(path string) (*Node, error) {
	// Short-circuit if the node already exists!
	if n, err := nt.FindNode(path); err == nil {
		if n.IsFolder() {
			return n, err
		}
		log.Errorf("%s: %s", constants.ErrFileExistsAndIsNotFolder, path)
		return nil, constants.ErrFileExistsAndIsNotFolder
	}

	// chop off the first PathSeparator.
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		path = path[1:]
	}
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) == 0 {
		log.Errorf("%s: %s", constants.ErrCannotCreateRootNode, path)
		return nil, constants.ErrCannotCreateRootNode
	}

	var (
		node     = nt.Node
		nextNode *Node
		err      error
	)

	for i, part := range parts {
		nextNode, err = nt.FindNode(strings.Join(parts[:i+1], string(os.PathSeparator)))
		if err != nil && err != constants.ErrNodeNotFound {
			return nil, err
		}
		if err == constants.ErrNodeNotFound {
			nextNode, err = node.CreateFolder(part)
			if err != nil {
				return nil, err
			}
		}

		if !nextNode.IsFolder() {
			log.Errorf("%s: %s", constants.ErrCannotCreateANodeUnderAFile, strings.Join(parts[:i+1], string(os.PathSeparator)))
			return nil, constants.ErrCannotCreateANodeUnderAFile
		}

		node = nextNode
	}

	return node, nil
}

func (nt *NodeTree) refresh() error {
	postURL := nt.client.GetMetadataURL("changes")
	c := &changes{
		Checkpoint: nt.Checkpoint,
	}
	// we are making a request to get the nt.Checkpoint
	if nt.Checkpoint == "" {
		c.Chunksize = 1
		c.MaxNodes = 1
	}
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		log.Errorf("error JSON encoding the new node metadata: %s", err)
		return constants.ErrJSONEncoding
	}

	// return format should be:
	// {"checkpoint": str, "reset": bool, "nodes": []}
	// {"checkpoint": str, "reset": false, "nodes": []}
	// {"end": true}
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("error creating the request: %s", err)
		return constants.ErrCreatingHTTPRequest
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := nt.client.Do(req)
	if err != nil {
		log.Errorf("error fetching the changes: %s", err)
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
			log.Errorf("error decoding the JSON body for the changes API: %s", err)
			return constants.ErrJSONDecodingResponseBody
		}
		nt.Checkpoint = cr.Checkpoint
		if cr.End {
			break
		}
		if cr.Reset {
			log.Debug("reset is required")
			return constants.ErrMustRefresh
		}
		if err := nt.updateNodes(cr.Nodes); err != nil {
			return err
		}
	}

	return nil
}

func (nt *NodeTree) updateNodes(nodes []*Node) error {
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
			for _, parentID := range newNode.Parents {
				parent, err := nt.FindByID(parentID)
				if err != nil {
					if err != constants.ErrNodeNotFound {
						log.Debugf("error trying to get the parent with ID %s: %s", parentID, err)
					}

					continue
				}
				parent.RemoveChild(node)
			}

			continue
		}

		// TODO: Handle change in parent IDs gracefully!
		sort.Strings(node.Parents)
		sort.Strings(newNode.Parents)
		if !reflect.DeepEqual(node.Parents, newNode.Parents) {
			log.Debugf("The parents of the node %s have changed. %v => %v", node.ID, node.Parents, newNode.Parents)
			return constants.ErrMustRefresh
		}

		// finally update the node itself
		(*node) = *newNode
	}

	return nil
}

func (nt *NodeTree) setClient(n *Node) {
	n.client = nt.client
	for _, node := range n.Nodes {
		nt.setClient(node)
	}
}
