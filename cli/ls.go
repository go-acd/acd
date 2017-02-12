package cli

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/acd.v0/node"

	"github.com/codegangsta/cli"
)

var (
	paths []string

	lsCommand = cli.Command{
		Name:         "ls",
		Usage:        "list directory contents",
		Description:  "ls list directory contents, multiple directories can be given. A directory must be prefixed by acd://",
		Action:       lsAction,
		BashComplete: lsBashComplete,
		Before:       lsBefore,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "long, l",
				Usage: "list in long format",
			},
		},
	}
)

func init() {
	registerCommand(lsCommand)
}

func lsAction(c *cli.Context) {
	for _, p := range paths {
		nodes, err := acdClient.List(p)
		if err != nil {
			log.Fatal(err)
		}

		if len(paths) > 1 {
			fmt.Printf("%s:\n", p)
		}

		if c.Bool("long") {
			lsLong(nodes)
		} else {
			lsShort(nodes)
		}

		if len(paths) > 1 {
			fmt.Println()
		}
	}
}

func lsBefore(c *cli.Context) error {
	for _, arg := range c.Args() {
		if !strings.HasPrefix(arg, "acd://") {
			continue
		}

		paths = append(paths, strings.TrimPrefix(arg, "acd://"))
	}
	if len(paths) == 0 {
		return fmt.Errorf("ls: at least one path prefixed by acd:// is required. Given: %v", c.Args())
	}

	return nil
}

func lsBashComplete(c *cli.Context) {
}

func lsLong(nodes node.Nodes) {
	for _, n := range nodes {
		if n.IsDir() {
			fmt.Print("d")
		} else {
			fmt.Print("-")
		}

		fmt.Printf("\t%d", n.Size())
		fmt.Printf("\t%s", n.ModTime())
		fmt.Printf("\t%s\n", n.Name)
	}
}

func lsShort(nodes node.Nodes) {
	sep := ""
	for _, n := range nodes {
		name := strings.Replace(n.Name, "\\", "\\\\", -1)
		name = strings.Replace(name, "\"", "\\\"", -1)
		fmt.Printf("%s\"%s\"", sep, name)
		sep = " "
	}
	fmt.Println("")
}
