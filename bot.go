package slackopher

import (
	"github.com/nlopes/slack"
	"gopkg.in/urfave/cli.v2"
	"github.com/Sirupsen/logrus"
)

func Bot(c *cli.Context) error {
	api := slack.New(c.GlobalString("slack_key"))
	api.SetDebug(c.GlobalBool("debug"))
	rtm := api.NewRTM()
	historian := NewHistorian("http://127.0.0.1:9200")

	channels, err := api.GetChannels(false)
	if err != nil {
		panic(err)
	}
	for _, channel := range channels {
		historian.JoinChannel(&channel)
	}

	responseChannel := make(chan Response)
	commandChannel := make(chan *Command)

	go rtm.ManageConnection()

	// Handle Bot Commands
	go Commander(commandChannel, responseChannel)

	// Handle Bot Response
	go Responder(api, responseChannel)

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
