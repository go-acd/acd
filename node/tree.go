package node

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

type (
	// Tree represents a node tree.
	Tree struct {
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
)

// RemoveNode removes this node from the server and from the NodeTree.
func (nt *Tree) RemoveNode(n *Node) error {
	if err := n.Remove(); err != nil {
		return err
	}

	for _, parentID := range n.Parents {
		parent, err := nt.FindByID(parentID)
		if err != nil {
			log.Debugf("parent ID %s not found", parentID)
			continue
		}
		parent.RemoveChild(n)
	}

	return nil
}

// NewTree returns the root node (the head of the tree).
func NewTree(c client, cacheFile string) (*Tree, error) {
	nt := &Tree{
		cacheFile: cacheFile,
		client:    c,
	}
	if err := nt.loadOrFetch(); err != nil {
		return nil, err
	}
	if err := nt.saveCache(); err != nil {
		return nil, err
	}

	return nt, nil
}

// Close finalizes the NodeTree
func (nt *Tree) Close() error {
	return nt.saveCache()
}

// MkdirAll creates a directory named path, along with any necessary parents,
// and returns the directory node and nil, or else returns an error. If path is
// already a directory, MkdirAll does nothing and returns the directory node
// and nil.
func (nt *Tree) MkdirAll(path string) (*Node, error) {
	var (
		err        error
		folderNode = nt.Node
		logLevel   = log.GetLevel()
		nextNode   *Node
		node       *Node
	)

	// Short-circuit if the node already exists!
	{
		log.SetLevel(log.DisableLogLevel)
		node, err = nt.FindNode(path)
		log.SetLevel(logLevel)
	}
	if err == nil {
		if node.IsDir() {
			return node, err
		}
		log.Errorf("%s: %s", constants.ErrFileExistsAndIsNotFolder, path)
		return nil, constants.ErrFileExistsAndIsNotFolder
	}

	// chop off the first /.
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		log.Errorf("%s: %s", constants.ErrCannotCreateRootNode, path)
		return nil, constants.ErrCannotCreateRootNode
	}

	for i, part := range parts {
		{
			log.SetLevel(log.DisableLogLevel)
			nextNode, err = nt.FindNode(strings.Join(parts[:i+1], "/"))
			log.SetLevel(logLevel)
		}
		if err != nil && err != constants.ErrNodeNotFound {
			return nil, err
		}
		if err == constants.ErrNodeNotFound {
			nextNode, err = folderNode.CreateFolder(part)
			if err != nil {
				return nil, err
			}
		}

		if !nextNode.IsDir() {
			log.Errorf("%s: %s", constants.ErrCannotCreateANodeUnderAFile, strings.Join(parts[:i+1], "/"))
			return nil, constants.ErrCannotCreateANodeUnderAFile
		}

		folderNode = nextNode
	}

	return folderNode, nil
}

func (nt *Tree) setClient(n *Node) {
	n.client = nt.client
	for _, node := range n.Nodes {
		nt.setClient(node)
	}
}

func (nt *Tree) buildNodeMap(current *Node) {
	if nt.Node == current {
		nt.nodeMap = make(map[string]*Node)
	}
	nt.nodeMap[current.ID] = current
	for _, node := range current.Nodes {
		nt.buildNodeMap(node)
	}
}

func (nt *Tree) loadOrFetch() error {
	var err error
	if err = nt.loadCache(); err != nil {
		log.Debug(err)
		if err = nt.fetchFresh(); err != nil {
			return err
		}
	}

	if err = nt.Sync(); err != nil {
		switch err {
		case constants.ErrMustFetchFresh:
			if err = nt.fetchFresh(); err != nil {
				return err
			}
			return nt.Sync()
		default:
			return err
		}
	}

	return nil
}

func (nt *Tree) fetchFresh() error {
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
			log.Errorf("%s: %s", constants.ErrCreatingHTTPRequest, err)
			return constants.ErrCreatingHTTPRequest
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := nt.client.Do(req)
		if err != nil {
			log.Errorf("%s: %s", constants.ErrDoingHTTPRequest, err)
			return constants.ErrDoingHTTPRequest
		}

		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&nl); err != nil {
			log.Errorf("%s: %s", constants.ErrJSONDecodingResponseBody, err)
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
		nt.setClient(node)
		nodeMap[node.ID] = node
	}

	for _, node := range nodeMap {
		if node.Name == "" && node.IsDir() && len(node.Parents) == 0 {
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
