package integrationtest

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func TestSimpleUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	needCleaning = true

	var (
		readmeFile       = "fixtures/README"
		remoteReadmeFile = remotePath(readmeFile)
	)

	// open the README file
	in, err := os.Open(readmeFile)
	if err != nil {
		t.Fatal(err)
	}
	inhash := md5.New()
	in.Seek(0, 0)
	io.Copy(inhash, in)
	inmd5 := hex.EncodeToString(inhash.Sum(nil))

	// test uploading
	c, err := newCachedClient(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}
	in.Seek(0, 0)
	if err := c.Upload(remoteReadmeFile, in); err != nil {
		t.Errorf("error uploading %s to %s: %s", readmeFile, remoteReadmeFile, err)
	}

	// test the NodeTree is updated
	node, err := c.NodeTree.FindNode(remoteReadmeFile)
	if err != nil {
		t.Errorf("c.NodeTree.FindNode(%q): got error %s", remoteReadmeFile, err)
	}
	if want, got := inmd5, node.ContentProperties.MD5; want != got {
		t.Errorf("c.NodeTree.FindNode(%q).ContentProperties.MD5: want %s got %s", remoteReadmeFile, want, got)
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
	node, err = c.NodeTree.FindNode(remoteReadmeFile)
	if err != nil {
		t.Errorf("reloaded cache, c.NodeTree.FindNode(%q): got error %s", remoteReadmeFile, err)
	}
	if want, got := inmd5, node.ContentProperties.MD5; want != got {
		t.Errorf("reloaded cache, c.NodeTree.FindNode(%q).ContentProperties.MD5: want %s got %s", remoteReadmeFile, want, got)
	}

	// check the file exists on the server
	uc, err := newUncachedClient()
	if err != nil {
		t.Fatal(err)
	}
	if err := uc.FetchNodeTree(); err != nil {
		t.Fatal(err)
	}
	out, err := c.Download(remoteReadmeFile)
	if err != nil {
		t.Errorf("error uploading %s to %s: %s", readmeFile, remoteReadmeFile, err)
	}
	outhash := md5.New()
	io.Copy(outhash, out)
	outmd5 := hex.EncodeToString(outhash.Sum(nil))

	if want, got := inmd5, outmd5; want != got {
		t.Errorf("c.Upload() hashes: want %s got %s", want, got)
	}
}
