package slack

import (
	"fmt"
	"github.com/bluele/slack"
)

// Alert Message
func AlertMessage(token string, channelName string, message string) {
	api := slack.New(token)
	channel, err := api.FindChannelByName(channelName)
	if err != nil {
		fmt.Println(err)
	}

	err = api.ChatPostMessage(channel.Id, message, nil)
	if err != nil {
		fmt.Println(err)
	}
}
