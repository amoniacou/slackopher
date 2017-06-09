package slackopher

import (
	"github.com/nlopes/slack"
	"strconv"
)

func LookupCommand(commandArray []string, response Response, responseChannel chan Response) {
	limit := 5
	// term := commandArray[2]
	if len(commandArray) > 3 {
		if intNumber, err := strconv.Atoi(commandArray[3]); err != nil {
			limit = intNumber
		}
	}
	fields := make([]slack.AttachmentField, limit)
	attachment := &slack.Attachment{
		Pretext: "Search result",
		Color: "#000000",
		Fields: fields,
	}
	response.Attachment = attachment
	responseChannel <- response
}
