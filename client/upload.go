package client

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// Upload uploads io.Reader to the path defined by the filename. It will create any non-existing
// folders.
func (c *Client) Upload(filename string, r io.Reader) error {
	if _, err := c.NodeTree.FindNode(filename); err == nil {
		return err
	}

	node, err := c.NodeTree.CreatePath(path.Dir(filename))
	if err != nil {
		return err
	}

	_, err = node.Upload(path.Base(filename), r)
	if err != nil {
		return err
	}

	return nil
}

// UploadFolder uploads an entire folder.
// If recursive is true, it will recurse through the entire filetree under
// localPath.
// If overwrite is false and an existing file with the same md5 was found, an
// error will be returned.
func (c *Client) UploadFolder(localPath, remotePath string, recursive, overwrite bool) error {
	log.Debugf("uploading %q to %q", localPath, remotePath)

	if err := filepath.Walk(localPath, c.uploadFolderFunc(localPath, remotePath, recursive, overwrite)); err != nil {
		return err
	}

	return nil
}

func (c *Client) uploadFolderFunc(localPath, remoteBasePath string, recursive, overwrite bool) filepath.WalkFunc {
	return func(fpath string, info os.FileInfo, err error) error {
		parts := strings.SplitAfter(fpath, localPath)
		remoteFilename := remoteBasePath + strings.Join(parts[1:], "/")
		remotePath := path.Dir(remoteFilename)
		log.Debugf("localPath %q remotePath %q fpath %q remoteFilename %q recursive %t overwrite %t",
			localPath, remotePath, fpath, remoteFilename, recursive, overwrite)

		// is this a folder?
		if info.IsDir() {
			log.Debugf("%q is a folder, skipping", fpath)
			return nil
		}
		// are we not recursive and trying to upload a file down the tree?
		if !recursive && localPath != path.Dir(fpath) {
			log.Debugf("%q is inside a sub-folder but we are not running recursively")
			return nil
		}

		log.Infof("uploading %q to %q", fpath, remoteFilename)
		remoteNode, err := c.NodeTree.CreatePath(remotePath)
		if err != nil {
			return err
		}

		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer f.Close()

		// does the file already exist?
		if fileNode, err := c.NodeTree.FindNode(remoteFilename); err == nil {
			if fileNode.IsFolder() {
				log.Errorf("%s: remoteFilename %q", constants.ErrFileExistsAndIsFolder, remoteFilename)
				return constants.ErrFileExistsAndIsFolder
			}
			hash := md5.New()
			f.Seek(0, 0)
			io.Copy(hash, f)
			if hex.EncodeToString(hash.Sum(nil)) == fileNode.ContentProperties.MD5 {
				log.Debugf("%q already exists and has the same content", fpath)
				return nil
			}

			log.Debugf("%q already exists, overwrite is %t", fpath, overwrite)
			if !overwrite {
				log.Errorf("%s: remoteFilename %q", constants.ErrFileExistsWithDifferentContents, remoteFilename)
				return constants.ErrFileExistsWithDifferentContents
			}

			f.Seek(0, 0)
			return fileNode.Overwrite(f)
		}

		f.Seek(0, 0)
		if _, err := remoteNode.Upload(path.Base(fpath), f); err != nil {
			return err
		}

		return nil
	}
}
