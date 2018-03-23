package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/meritlabs/overlord/overlord/messaging"
	"github.com/spf13/viper"
)

var (
	pingDelay  = 30
	checkDelay = 90
	ticks      = checkDelay / pingDelay
)

func main() {
	fmt.Printf("Overlord is ready.\n")

	var configPath string

	flag.StringVar(&configPath, "config", "/etc", "Daemon mode. Can be mainnet, testnet or regtest")
	flag.Parse()

	viper.SetConfigName("overlord")
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	slackAPIKey := viper.GetString("slack_api_key")
	slackChannel := viper.GetString("slack_channel")
	ips := viper.GetStringSlice("ip_addresses")

	fmt.Printf("Monitoring IPs: %s \n", strings.Join(ips, ", "))
	fmt.Printf("Posting to Slack channel: %s \n", slackChannel)

	slackAPI := messaging.InitSlack(slackAPIKey)

	countdown := ticks
	for t := range time.NewTicker(time.Duration(pingDelay) * time.Second).C {
		fmt.Printf("Performing checks %v, %v\n", t, countdown)
		countdown--

		go messaging.DoPing(ips)

		if countdown == 0 {
			go messaging.DoCheck(slackAPI, slackChannel, ips)
			countdown = ticks
		}
	}
}
