package main

import (
	"gopkg.in/urfave/cli.v2"
	"log"
	"github.com/nlopes/slack"
	"github.com/Sirupsen/logrus"
	"github.com/amoniacou/slackopher"
)

var botCmd = cli.Command{
	Name: "bot",
	Usage: "starts slaclopher bot daemon",
	Aliases: []string{"b"},
	Action: func(c *cli.Context) {
		if err := bot(c); err != nil {
			log.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "slack_key",
			EnvVar: "SLACK_KEY",
			Usage: "Slack API key",
		},
	},
}

func bot(c *cli.Context) error {
	api := slack.New(c.String("slack_key"))
	api.SetDebug(c.GlobalBool("debug"))
	rtm := api.NewRTM()
	historian := slackopher.NewHistorian(c.GlobalString("elastic_url"))

	channels, err := api.GetChannels(false)
	if err != nil {
		panic(err)
	}
	for _, channel := range channels {
		historian.JoinChannel(&channel)
	}

	responseChannel := make(chan slackopher.Response)
	commandChannel := make(chan *slackopher.Command)

	go rtm.ManageConnection()

	// Handle Bot Commands
	go slackopher.Commander(commandChannel, responseChannel)

	// Handle Bot Response
	go slackopher.Responder(api, responseChannel)

	Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev:= msg.Data.(type) {
			case *slack.ConnectedEvent:
				historian.BotId = ev.Info.User.ID
				historian.BotName = ev.Info.User.Name
				logrus.Info("Connected: " + ev.Info.User.Name)
			case *slack.MessageEvent:
				logrus.Info("Process message")
				historian.ProcessMessage(api, ev, commandChannel)
			case *slack.RTMError:
			case *slack.InvalidAuthEvent:
				break Loop
			default:

			}
		}
	}
	return nil
}
