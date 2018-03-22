package messaging

import (
	"fmt"

	"github.com/nlopes/slack"
)

var (
	slackApi *slack.Client
)

func InitSlack(token string) {
	slackApi := slack.New(token)
	fmt.Printf("Slack API initialized: %v \n", slackApi)
}
