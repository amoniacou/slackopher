package main

import (
	"gopkg.in/urfave/cli.v2" // imports as package "cli"
	"github.com/amoniacou/slackopher"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "slackopher"
	app.Version = "0.0.1"
	app.Usage = "Slackopher - Slack history bot"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEBUG",
			Name:   "debug",
			Usage:  "start app in debug mode",
		},
		cli.StringFlag{
			Name: "slack_key",
			EnvVar: "SLACK_KEY",
			Usage: "Slack API key",
		},
	}
	app.Action = slackopher.Bot
	app.Run(os.Args)
}
