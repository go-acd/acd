package acd

import (
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// Download returns an io.ReadCloser for path. The caller is responsible for
// closing the body.
func (c *Client) Download(path string) (io.ReadCloser, error) {
	log.Debugf("downloading %q", path)

	node, err := c.NodeTree.FindNode(path)
	if err != nil {
		return nil, err
	}

	return node.Download()
}

// DownloadFolder downloads an entire folder to a path, if recursive is true,
// it will also download all subfolders.
func (c *Client) DownloadFolder(localPath, remotePath string, recursive bool) error {
	log.Debugf("downloading %q to %q", localPath, remotePath)

	if err := os.Mkdir(localPath, os.FileMode(0755)); err != nil && !os.IsExist(err) {
		log.Errorf("%s: %s", constants.ErrCreateFolder, err)
		return constants.ErrCreateFolder
	}
	rootNode, err := c.GetNodeTree().FindNode(remotePath)
	if err != nil {
		return nil
	}
	for _, node := range rootNode.Nodes {
		flp := path.Join(localPath, node.Name)
		frp := fmt.Sprintf("%s/%s", remotePath, node.Name)
		if node.IsDir() {
			if recursive {
				if err := c.DownloadFolder(flp, frp, recursive); err != nil {
					return err
				}
			}

			continue
		}

		con, err := node.Download()
		if err != nil {
			return err
		}
		f, err := os.Create(flp)
		if err != nil {
			log.Errorf("%s: %s", constants.ErrCreateFile, flp)
			return constants.ErrCreateFile
		}
		log.Debugf("saving %s as %s", frp, flp)
		_, err = io.Copy(f, con)
		f.Close()
		con.Close()
		if err != nil {
			log.Errorf("%s: %s", constants.ErrWritingFileContents, err)
			return err
		}
	}

	return nil
}
