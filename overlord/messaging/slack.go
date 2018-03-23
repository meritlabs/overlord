package messaging

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

func InitSlack(token string) *slack.Client {
	slackAPI := slack.New(token)
	return slackAPI
}

func ReadChannels(slackAPI *slack.Client) {
	channels, err := slackAPI.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Println(channel.Name)
		// channel is of type conversation & groupConversation
		// see all available methods in `conversation.go`
	}
}

func PostMessage(slackAPI *slack.Client, message string, channel string) {
	params := slack.PostMessageParameters{}
	channelID, timestamp, err := slackAPI.PostMessage(channel, message, params) //grab channel from config

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func PostOddNodesReport(slackAPI *slack.Client, channel string, nodes OddNodes) {
	var message string
	post := false

	if len(nodes.RespondedWithError) > 0 {
		message += fmt.Sprintf("Unhealthy or misconfigured nodes: %v \n", strings.Join(nodes.RespondedWithError, ", "))
		post = true
	}

	if len(nodes.HeadersAndBlocksDontMatch) > 0 {
		message += fmt.Sprintf("Nodes that have different number of blocks and headers: %v \n", strings.Join(nodes.HeadersAndBlocksDontMatch, ", "))
		post = true
	}

	if len(nodes.HaveDifferentNumberOfHeaders) > 0 {
		message += fmt.Sprintf("Nodes that have different number of headers compared to other nodes:\n")
		for key, nodes := range nodes.HaveDifferentNumberOfHeaders {
			message += fmt.Sprintf("- %d: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if len(nodes.HaveDifferentNumberOfBlocks) > 0 {
		message += fmt.Sprintf("Nodes that have different number of blocks compared to other nodes:\n")
		for key, nodes := range nodes.HaveDifferentNumberOfBlocks {
			message += fmt.Sprintf("- %d: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if len(nodes.HaveDifferentDifficulty) > 0 {
		message += fmt.Sprintf("Nodes that have different difficulty compared to other nodes:\n")
		for key, nodes := range nodes.HaveDifferentDifficulty {
			message += fmt.Sprintf("- %f: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if len(nodes.HaveDifferentChainwork) > 0 {
		message += fmt.Sprintf("Nodes that have different chainwork compared to other nodes:\n")
		for key, nodes := range nodes.HaveDifferentChainwork {
			message += fmt.Sprintf("- %s: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if len(nodes.HaveDifferentBestBlockHash) > 0 {
		message += fmt.Sprintf("Nodes that have different best block hash compared to other nodes:\n")
		for key, nodes := range nodes.HaveDifferentBestBlockHash {
			message += fmt.Sprintf("- %s: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if post {
		PostMessage(slackAPI, message, channel)
	}
}
