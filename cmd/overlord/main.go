package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	controllers "github.com/meritlabs/overlord/pkg/controllers/statuspage"
	"github.com/meritlabs/overlord/pkg/messaging"
	"github.com/spf13/viper"
)

const (
	pingDelay          = 30
	checkDelay         = 600
	ticks              = checkDelay / pingDelay
	ticksToTestVersion = 10
)

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range controllers.Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func startHTTPServer() {
	fmt.Println("Running HTTP Server")

	r := gin.Default()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/html/index.tmpl", gin.H{})
	})

	r.Run(":9999")
}

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

	startHTTPServer()

	countdown := ticks
	versionCheckCounter := 0
	for t := range time.NewTicker(time.Duration(pingDelay) * time.Second).C {
		fmt.Printf("Performing checks %v, %v\n", t, countdown)
		countdown--

		go messaging.DoPing(ips)

		if countdown == 0 {
			go messaging.DoCheck(slackAPI, slackChannel, ips)
			countdown = ticks
			versionCheckCounter++
		}

		if versionCheckCounter == ticksToTestVersion {
			go messaging.DoVersionCheck(slackAPI, slackChannel, ips)
			versionCheckCounter = 0
		}
	}
}
