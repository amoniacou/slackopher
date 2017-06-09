package slackopher

import (
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"time"
	"errors"
)

type Channel struct {
	Info *slack.Channel
}

type Channels struct {
	Cache *cache.Cache
}

func NewChannels() *Channels {
	return &Channels{
		Cache:cache.New(600 * time.Minute, 1200 * time.Minute),
	}
}

func (c *Channels) GetChannel(id string) (Channel, error) {
	channel, found := c.Cache.Get(id)
	if found {
		return channel.(Channel), nil
	} else {
		return Channel{}, errors.New("Channel not found")
	}
}

func (c *Channels) AddChannel(channel *slack.Channel) Channel {
	ch := Channel{
		Info: channel,
	}
	c.Cache.Set(channel.ID, ch, cache.NoExpiration)
	return ch
}
