package integrationtest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSimpleFolderUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	needCleaning = true

	testUploadFolder(t, "fixtures/simplefolder", false, false)
}

func TestRecursiveFolderUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	needCleaning = true

	testUploadFolder(t, "fixtures/recursivefolder", true, false)
}

func TestConflictFolderUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	needCleaning = true

	testUploadFolder(t, "fixtures/conflictfolder/v1", true, false)
	testUploadFolder(t, "fixtures/conflictfolder/v2", true, true)
}

func testUploadFolder(t *testing.T, localFolder string, recursive, overwrite bool) {
	var (
		remoteFolder = remotePath(localFolder)
		files        = listFiles(localFolder)
		md5s         = make(map[string]string, len(files))
	)

	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s/%s", localFolder, file))
		if err != nil {
			t.Fatal(err)
		}
		hash := md5.New()
		io.Copy(hash, f)
		md5s[file] = hex.EncodeToString(hash.Sum(nil))
	}

	// test uploading
	c, err := newCachedClient(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}
	if err := c.UploadFolder(localFolder, remoteFolder, recursive, overwrite); err != nil {
		t.Errorf("error uploading %s to %s: %s", localFolder, remoteFolder, err)
	}

	// test the NodeTree is updated
	for _, file := range files {
		remoteFile := fmt.Sprintf("%s/%s", remoteFolder, file)
		node, err := c.NodeTree.FindNode(remoteFile)
		if err != nil {
			t.Errorf("c.NodeTree.FindNode(%s): got error %s", remoteFile, err)
		}
		// run the following test only if we find the node
		if err == nil {
			if want, got := md5s[file], node.ContentProperties.MD5; want != got {
				t.Errorf("c.NodeTree.FindNode(%s).ContentProperties.MD5: want %s got %s", file, want, got)
			}
		}
	}

	// test the cache is being saved updated
	c.Close()
	c, err = newCachedClient(false)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		remoteFile := fmt.Sprintf("%s/%s", remoteFolder, file)
		node, err := c.NodeTree.FindNode(remoteFile)
		if err != nil {
			t.Errorf("reloaded cache, c.NodeTree.FindNode(%s): got error %s", file, err)
		}
		// run the following test only if we find the node
		if err == nil {
			if want, got := md5s[file], node.ContentProperties.MD5; want != got {
				t.Errorf("reloaded cache, c.NodeTree.FindNode(%s).ContentProperties.MD5: want %s got %s", file, want, got)
			}
		}
	}

	// download the folder from the server and verify the contents byte-per-byte
	localPath, err := ioutil.TempDir("", "acd-folder-upload-test-")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.DownloadFolder(localPath, remoteFolder, recursive); err != nil {
		t.Errorf("c.DownloadFolder(%s, %s, %t) error: %s", localPath, remoteFolder, recursive, err)
	}
	for _, file := range files {
		localFile := path.Join(localPath, file)
		f, err := os.Open(localFile)
		if err != nil {
			if os.IsNotExist(err) {
				t.Errorf("cannot open the downloaded file at %s: does not exist", localFile)
			} else {
				t.Errorf("cannot open the downloaded file at %s: %s", localFile, err)
			}
		}
		hash := md5.New()
		io.Copy(hash, f)

		if want, got := md5s[file], hex.EncodeToString(hash.Sum(nil)); want != got {
			t.Errorf("c.DownloadFolder(%s, %s, %t), md5 of %s: want %s got %s", localPath, remoteFolder, recursive, localFile, want, got)
		}
	}
	if err := os.RemoveAll(localPath); err != nil {
		t.Logf("error removing the temporary folder %q, please remove it manually: %s", localPath, err)
	}
}
