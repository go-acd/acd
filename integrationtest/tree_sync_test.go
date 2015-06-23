package integrationtest

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func TestTreeSync(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	needCleaning = true

	var (
		readmeFile       = "fixtures/syncfolder/README.syncfolder"
		remoteReadmeFile = remotePath(readmeFile)
	)

	// create two cached clients
	c1, err := newCachedClient(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := c1.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}
	c2, err := newCachedClient(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := c2.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}

	// using the first client upload the README file
	in, err := os.Open(readmeFile)
	if err != nil {
		t.Fatal(err)
	}
	inhash := md5.New()
	in.Seek(0, 0)
	io.Copy(inhash, in)
	inmd5 := hex.EncodeToString(inhash.Sum(nil))
	in.Seek(0, 0)
	if err := c1.Upload(remoteReadmeFile, in); err != nil {
		t.Errorf("error uploading %s to %s: %s", readmeFile, remoteReadmeFile, err)
	}

	// using the second client, sync and find the node
	if err := c2.NodeTree.Sync(); err != nil {
		t.Fatal(err)
	}
	readmeNode, err := c2.NodeTree.FindNode(remoteReadmeFile)
	if err != nil {
		t.Fatalf("c2.NodeTree.FindNode(%q) error: %s", remoteReadmeFile, err)
	}
	if want, got := inmd5, readmeNode.ContentProperties.MD5; want != got {
		t.Errorf("c.NodeTree.FindNode(%q).ContentProperties.MD5: want %s got %s", remoteReadmeFile, want, got)
	}
}
