package slackopher

type Message struct {
	User        string `json:"user"`
	UserId      string `json:"user_id"`
	Avatar      string `json:"avatar"`
	Channel     string `json:"channel"`
	ChannelName string `json:"channel_name"`
	Body        string `json:"body"`
	Timestamp   string `json:"ts"`
}
