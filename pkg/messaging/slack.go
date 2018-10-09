package messaging

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

// InitSlack - initiate Slack API
func InitSlack(token string) *slack.Client {
	slackAPI := slack.New(token)
	return slackAPI
}

// ReadChannels - get Slack channels list
func ReadChannels(slackAPI *slack.Client) {
	channels, err := slackAPI.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Println(channel.Name)
	}
}

// PostMessage - sends a message to Slack channel
func PostMessage(slackAPI *slack.Client, message string, channel string) {
	params := slack.PostMessageParameters{}
	channelID, timestamp, err := slackAPI.PostMessage(channel, message, params) //grab channel from config

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

// PostOddNodesReport - posts report containing all nodes with differences in blockchain status
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

// PostVersionMismatchReport - posts report containing all nodes with differences daemon or protocol version
func PostVersionMismatchReport(slackAPI *slack.Client, channel string, nodes MismatchVeresionNodes) {
	var message string
	post := false

	if len(nodes.RespondedWithError) > 0 {
		message += fmt.Sprintf("Unhealthy or misconfigured nodes: %v \n", strings.Join(nodes.RespondedWithError, ", "))
		post = true
	}

	if len(nodes.HaveDifferentVersions) > 0 {
		message += fmt.Sprintf("Newest Version is: %v \n", nodes.NewestVersion)
		message += fmt.Sprintf("Nodes that have different version compared to the newest nodes:\n")
		for key, nodes := range nodes.HaveDifferentVersions {
			message += fmt.Sprintf("- %d: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if len(nodes.HaveDifferentProtocolVersions) > 0 {
		message += fmt.Sprintf("Newest Protocol Version is: %v \n", nodes.NewestProtocolVersion)
		message += fmt.Sprintf("Nodes that have different protocol version compared to the newest nodes:\n")
		for key, nodes := range nodes.HaveDifferentProtocolVersions {
			message += fmt.Sprintf("- %d: %v \n", key, strings.Join(nodes, ", "))
		}
		post = true
	}

	if post {
		PostMessage(slackAPI, message, channel)
	}
}
