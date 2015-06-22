package integrationtest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"gopkg.in/acd.v0"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

const (
	devNullCacheFile   string = "/dev/null"
	testFolderBasePath string = "/acd_test_folder"
	baseTokenFile      string = "acd-token.json"
)

var (
	cacheFile      string
	cacheFiles     []string
	configFiles    []string
	tokenFiles     []string
	testFolderPath string
	needCleaning   bool
)

func TestMain(m *testing.M) {
	defer func() {
		if r := recover(); r != nil {
			cleanUp()
		}
	}()

	cacheFile = newTempFile("acd-cache-")
	cacheFiles = append(cacheFiles, cacheFile)
	testFolderPath = fmt.Sprintf("%s/%d", testFolderBasePath, time.Now().UnixNano())

	// disable all logs
	log.SetLevel(log.ErrorLevel)

	// run all the tests
	code := m.Run()

	// Cleanup after the run
	cleanUp()

	// exit with the return status
	os.Exit(code)
}

func cleanUp() {
	if needCleaning {
		// remove the test folder
		if err := removeTestFolder(); err != nil {
			log.Errorf("error removing the test folder: %s", err)
		}

		// avoid double cleaning
		needCleaning = false
	}

	// remove all cache files
	for _, cf := range cacheFiles {
		os.Remove(cf)
	}

	// remove all config files.
	for _, cf := range configFiles {
		os.Remove(cf)
	}

	// remove all token files.
	for _, cf := range tokenFiles {
		os.Remove(cf)
	}
}

func newTempFile(baseName string) string {
	f, _ := ioutil.TempFile("", baseName)
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func newCachedClient(ncf bool) (*acd.Client, error) {
	if ncf {
		cacheFile = newTempFile("acd-cache-")
		cacheFiles = append(cacheFiles, cacheFile)
	}
	return acd.New(newConfigFile(cacheFile))
}

func newUncachedClient() (*acd.Client, error) {
	return acd.New(newConfigFile(devNullCacheFile))
}

func newConfigFile(cacheFile string) string {
	tokenFile := newTempFile("acd-token-")
	tokenFiles = append(tokenFiles, tokenFile)
	configFile := newTempFile("acd-config-")
	configFiles = append(configFiles, configFile)

	of, err := os.Open(baseTokenFile)
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()
	nf, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer nf.Close()
	io.Copy(nf, of)

	cf, err := os.Create(configFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &acd.Config{
		TokenFile: tokenFile,
		CacheFile: cacheFile,
	}
	defer cf.Close()
	if err := json.NewEncoder(cf).Encode(config); err != nil {
		log.Fatal(err)
	}

	return configFile
}

func removeTestFolder() error {
	c, err := newUncachedClient()
	if err != nil {
		return err
	}
	if err := c.FetchNodeTree(); err != nil {
		return err
	}
	node, err := c.NodeTree.FindNode(testFolderPath)
	if err != nil && err != constants.ErrNodeNotFound {
		return err
	}
	if node == nil {
		return constants.ErrNodeNotFound
	}
	if node.Name != path.Base(testFolderPath) {
		return fmt.Errorf("something is wrong, the node's name is %s and not %s", node.Name, testFolderPath)
	}

	return c.NodeTree.RemoveNode(node)
}

func remotePath(fp string) string {
	p := strings.Replace(fp, "fixtures/", "", 1)
	r, err := regexp.Compile(`/(v1|v2)`)
	if err != nil {
		panic(err)
	}
	if ok := r.MatchString(p); ok {
		p = strings.Replace(p, "/v1", "/", 1)
		p = strings.Replace(p, "/v2", "/", 1)
	}
	p = strings.TrimSuffix(p, "/")
	return fmt.Sprintf("%s/%s", testFolderPath, p)
}

// listFiles returns the list of all of the files in folder and it's subfolders
// but it does not include the subfolders as entries.
func listFiles(folder string) []string {
	var files []string
	filepath.Walk(folder, func(fp string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		parts := strings.SplitAfter(fp, fmt.Sprintf("%s%c", folder, os.PathSeparator))
		nfp := strings.Join(parts[1:], string(os.PathSeparator))
		files = append(files, nfp)
		return nil
	})

	return files
}
