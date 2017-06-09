package slackopher

import (
	elastic "gopkg.in/olivere/elastic.v5"
	"github.com/nlopes/slack"
	"context"
	"strings"
	"github.com/Sirupsen/logrus"
	"log"
)

type Historian struct {
	ElasticClient *elastic.Client
	Ctx           context.Context
	BotId         string
	BotName       string
	Users         *Users
	Channels      *Channels
}

func NewHistorian(elastic_url string) *Historian {
	client, err := elastic.NewClient(elastic.SetURL(elastic_url))
	if err != nil {
		panic(err)
	}
	historian := &Historian{
		ElasticClient: client,
		Ctx: context.Background(),
	}
	historian.Setup()
	return historian
}

func (h *Historian) Setup() error {
	exists, err := h.ElasticClient.IndexExists("messages").Do(h.Ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		h.CreateIndex()
	}
	h.Users = NewUsers()
	h.Channels = NewChannels()
	return nil
}

func (h *Historian) CreateIndex() {
	mapping := `{
		"mappings":{
			"message":{
				"properties":{
					"user":{
						"type":"string"
					},
					"channel":{
						"type":"string"
					},
					"body": {
						"type":"text"
					},
					"ts": {
						"type":"date",
						"format": "strict_date_optional_time||epoch_millis"
					}
				}
			}
		}
	}`
	_, err := h.ElasticClient.CreateIndex("messages").BodyString(mapping).Do(h.Ctx)
	if err != nil {
		panic(err)
	}
}

func (h *Historian) JoinChannel(channel *slack.Channel) {
	h.Channels.AddChannel(channel)
}

func (h *Historian) SendCommand(msg *slack.MessageEvent, channel chan *Command) {
	logrus.Info("Send New Command")
	channel <- &Command{
		ChannelId: msg.Channel,
		Event: msg,
		UserId: msg.User,
	}
}

func (h *Historian) SaveMessage(api *slack.Client, msg *slack.MessageEvent) {
	logrus.Info("News Message")
	u, err := h.Users.GetUser(msg.User)
	if err != nil {
		user, err := api.GetUserInfo(msg.User)
		if err != nil {
			log.Fatal("Could not retrieve user")
			return
		}
		u = h.Users.AddUser(user)
	}
	channel, err := h.Channels.GetChannel(msg.Channel)
	if err != nil {
		log.Fatal("Could not retreive channel")
		return
	}
	message := Message{
		Body: msg.Text,
		UserId: msg.User,
		User: u.Info.Name,
		Avatar: u.Info.Profile.Image48,
		Channel: msg.Channel,
		ChannelName: channel.Info.Name,
		Timestamp: strings.Split(msg.Timestamp, ".")[0],
	}
	logrus.Info(message)
	logrus.Info(message.Body)
	_, err = h.ElasticClient.Index().
		Index("messages").
		Type("message").
		BodyJson(message).
		Do(h.Ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Historian) ProcessMessage(api *slack.Client, msg *slack.MessageEvent, channel chan *Command) {
	if msg.Type != "message" {
		return
	}
	if strings.HasPrefix(msg.Text, "<@" + h.BotId + ">") {
		h.SendCommand(msg, channel)
	} else {
		if msg.SubType == "channel_join" && msg.User == h.BotId {
			logrus.Info("Join new channel")
			channelInfo, err := api.GetChannelInfo(msg.Channel)
			if err != nil {
				log.Fatal(err)
			} else {
				h.JoinChannel(channelInfo)
			}
		}
		h.SaveMessage(api, msg)
	}
}
