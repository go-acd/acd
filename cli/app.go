package cli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"gopkg.in/acd.v0"
	"gopkg.in/acd.v0/internal/log"
)

// TODO(kalbasit): I do not like this API code not even a bit.
// a) I used codegangsta/cli wrong or overthought it.
// b) codegangsta/cli is not the right library for this project.
// This entire package should be re-written and TESTED!.

var (
	commands  []cli.Command
	acdClient *acd.Client
)

// New creates a new CLI application.
func New() *cli.App {
	app := cli.NewApp()
	app.Author = "Wael Nasreddine"
	app.Email = "wael.nasreddine@gmail.com"
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Name = "acd"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, c",
			Value: acd.DefaultConfigFile(),
			Usage: "the path of the configuration file",
		},

		cli.IntFlag{
			Name:  "log-level, l",
			Value: int(log.FatalLevel),
			Usage: fmt.Sprintf("possible log levels: %s", log.Levels()),
		},
	}

	app.Before = beforeCommand
	app.After = afterCommand
	app.Commands = commands
	return app
}

func registerCommand(c cli.Command) {
	commands = append(commands, c)
}

func beforeCommand(c *cli.Context) error {
	var err error

	// set the log level
	log.SetLevel(log.Level(c.Int("log-level")))

	// create a new client
	if acdClient, err = acd.New(c.String("config-file")); err != nil {
		return fmt.Errorf("error creating a new ACD client: %s", err)
	}

	// fetch the nodetree
	if err = acdClient.FetchNodeTree(); err != nil {
		return fmt.Errorf("error fetch the node tree: %s", err)
	}

	return nil
}

func afterCommand(_ *cli.Context) error {
	if acdClient != nil {
		return acdClient.Close()
	}
	return nil
}
