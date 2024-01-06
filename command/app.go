package command

import (
	"github.com/kkoch986/ai-skeletons-playback/command/internal/server"
	cli "github.com/urfave/cli/v2"
)

// App presents the command line interface to start the different
// processes.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "playback"
	app.Description = "a service for playing audio in sync with sequences generated elsewhere in the pipeline"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"LOG_LEVEL"},
			Usage:   "Set the log level",
			Value:   "WARN",
		},
	}
	app.Commands = []*cli.Command{
		server.Command,
	}
	return app
}
