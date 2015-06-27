package cli

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"

	"github.com/codegangsta/cli"
)

var (
	cpCommand = cli.Command{
		Name:         "cp",
		Usage:        "copy files",
		Description:  "cp copy files, multiple files can be given. It follows the usage of cp whereas the last entry is the destination and has to be a folder if multiple files were given",
		Action:       cpAction,
		BashComplete: cpBashComplete,
		Before:       cpBefore,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "recursive, R",
				Usage: "cp recursively",
			},
		},
	}

	action string
)

func init() {
	registerCommand(cpCommand)
}

func cpAction(c *cli.Context) {
	if strings.HasPrefix(c.Args()[len(c.Args())-1], "acd://") {
		cpUpload(c)
	} else {
		cpDownload(c)
	}
}

func cpUpload(c *cli.Context) {
	// make sure the destination is a folder if it exists upstream and more than
	// one file is scheduled to be copied.
	dest := strings.TrimPrefix(c.Args()[len(c.Args())-1], "acd://")
	destNode, err := acdClient.NodeTree.FindNode(dest)
	if err == nil {
		// make sure if the remote node exists, it is a folder.
		if len(c.Args()) > 2 {
			if !destNode.IsDir() {
				log.Fatalf("cp: target %q is not a directory", dest)
			}
		}
	}

	for _, src := range c.Args()[:len(c.Args())-1] {
		if strings.HasPrefix(src, "acd://") {
			fmt.Printf("cp: target %q is amazon, src cannot be amazon when destination is amazon. Skipping\n", src)
			continue
		}
		stat, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("cp: %s: %s", constants.ErrFileNotFound, src)
			}

			log.Fatalf("cp: %s: %s", constants.ErrStatFile, src)
		}
		if stat.IsDir() {
			if !c.Bool("recursive") {
				fmt.Printf("cp: %q is a directory (not copied).", src)
				continue
			}
			destFile := dest
			if destNode != nil {
				if !destNode.IsDir() {
					log.Fatalf("cp: target %q is not a directory", dest)
				}
				destFile = fmt.Sprintf("%s/%s", dest, path.Base(src))
			}
			acdClient.UploadFolder(src, destFile, true, true)
			continue
		}
		f, err := os.Open(src)
		if err != nil {
			log.Fatalf("%s: %s -- %s", constants.ErrOpenFile, err, src)
		}
		err = acdClient.Upload(dest, true, f)
		f.Close()
		if err != nil {
			log.Fatalf("%s: %s", err, dest)
		}
	}
}

func cpDownload(c *cli.Context) {
	dest := c.Args()[len(c.Args())-1]
	destDir := false
	destStat, err := os.Stat(dest)
	if err == nil && destStat.IsDir() {
		destDir = true
	}
	if len(c.Args()) > 2 {
		if err == nil && !destDir {
			log.Fatalf("cp: target %q is not a directory", dest)
		}
	}

	for _, src := range c.Args()[:len(c.Args())-1] {
		if !strings.HasPrefix(src, "acd://") {
			fmt.Printf("cp: source %q is local, src cannot be local when destination is local. Skipping\n", src)
			continue
		}
		srcPath := strings.TrimPrefix(src, "acd://")
		destPath := dest
		if destDir {
			destPath = path.Join(destPath, path.Base(srcPath))
		}
		srcNode, err := acdClient.GetNodeTree().FindNode(srcPath)
		if err != nil {
			fmt.Printf("cp: source %q not found. Skipping", src)
			continue
		}
		if srcNode.IsDir() {
			acdClient.DownloadFolder(destPath, srcPath, c.Bool("recursive"))
		} else {
			content, err := acdClient.Download(srcPath)
			if err != nil {
				fmt.Printf("cp: error downloading source %q. Skipping", src)
			}
			// TODO: respect umask
			if err := os.MkdirAll(path.Dir(destPath), os.FileMode(0755)); err != nil {
				fmt.Printf("cp: error creating the parents folders of %q: %s. Skipping", destPath, err)
				continue
			}
			// TODO: respect umask
			f, err := os.Create(destPath)
			if err != nil {
				fmt.Printf("cp: error writing %q: %s. Skipping", destPath, err)
				continue
			}
			io.Copy(f, content)
			f.Close()
		}
	}
}

func cpBashComplete(c *cli.Context) {
}

func cpBefore(c *cli.Context) error {
	foundRemote := false
	for _, arg := range c.Args() {
		if strings.HasPrefix(arg, "acd://") {
			foundRemote = true
			break
		}
	}

	if !foundRemote {
		return fmt.Errorf("cp: at least one path prefixed by acd:// is required. Given: %v", c.Args())
	}

	return nil
}
