package slackopher

import "github.com/nlopes/slack"

func HelpCommand(response Response, responseChannel chan Response) {
	fields := make([]slack.AttachmentField, 0)
	for k, v := range Commands {
		fields = append(fields, slack.AttachmentField{
			Title: "<bot> " + k,
			Value: v,
		})
		attachment := &slack.Attachment{
			Pretext: "Supported Command List",
			Color: "#000000",
			Fields: fields,
		}
		response.Attachment = attachment
		responseChannel <- response
	}
}
