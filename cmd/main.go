package main

import (
	"gopkg.in/urfave/cli.v2"
	"os"
)

var (
	Version string
	BuildTime string
)

func main() {
	app := cli.NewApp()
	app.Name = "slackopher"
	app.Version = Version
	app.Usage = "Slackopher - Slack history bot"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEBUG",
			Name:   "debug, D",
			Usage:  "start app in debug mode",
		},
		cli.StringFlag{
			Name: "elastic_url",
			EnvVar: "ELASTIC_URL",
			Value: "http://localhost:9200",
			Usage: "ElasticSearch URL",
		},
	}
	app.Commands = []cli.Command{
		webCmd,
		botCmd,
	}
	app.Run(os.Args)
}
