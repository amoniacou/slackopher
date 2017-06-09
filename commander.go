package slackopher

import (
	"github.com/nlopes/slack"
	"strings"
)

type Command struct {
	ChannelId string
	Event     *slack.MessageEvent
	UserId    string
}

var Commands = map[string]string{
	"help":"See the available bot commands.",
}

func Commander(commandChannel chan *Command, responseChannel chan Response) {
	var response Response
	for {
		command := <-commandChannel
		response.DisplayTitle = "History Lookup"
		response.ChannelId = command.ChannelId
		commandArray := strings.Fields(command.Event.Text)
		if len(commandArray) > 1 {
			switch commandArray[1] {
			case "help":
				HelpCommand(response, responseChannel)
			case "lookup":
				LookupCommand(commandArray, response, responseChannel)
			}
		} else {
			HelpCommand(response, responseChannel)
		}
	}
}
