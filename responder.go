package slackopher

import (
	"github.com/nlopes/slack"
	"log"
)

type Response struct {
	ChannelId    string
	Attachment   *slack.Attachment
	DisplayTitle string
}

func Responder(api *slack.Client,responseChannel chan Response) {
	for {
		response := <-responseChannel
		params := slack.PostMessageParameters{
			AsUser:true,
			Attachments: []slack.Attachment{*response.Attachment},
		}
		_, _, err := api.PostMessage(response.ChannelId, response.DisplayTitle, params)
		if err != nil {
			log.Fatal(err)
		}
	}
}