package messaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meritlabs/overlord/pkg/blockchain"
	"github.com/nlopes/slack"
)

type Ping struct {
	ip       string
	isOnline bool
}

type Check struct {
	ip       string
	response *blockchain.CheckResponse
	isOnline bool
}

type VersionCheck struct {
	ip       string
	response *blockchain.VersionResponse
	isOnline bool
}

type OddNodes struct {
	RespondedWithError           []string
	HeadersAndBlocksDontMatch    []string
	HaveDifferentNumberOfHeaders map[int32][]string
	HaveDifferentNumberOfBlocks  map[int32][]string
	HaveDifferentDifficulty      map[float64][]string
	HaveDifferentChainwork       map[string][]string
	HaveDifferentBestBlockHash   map[string][]string
}

type MismatchVeresionNodes struct {
	RespondedWithError            []string
	HaveDifferentVersions         map[int64][]string
	HaveDifferentProtocolVersions map[int64][]string
	NewestVersion                 int64
	NewestProtocolVersion         int64
}

func DoPing(ips []string) {
	pingChannel := make(chan Ping)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s:8080/ping", ip)
			res, err := http.Get(host)

			if err != nil {
				fmt.Printf("Error on ping request: %s \n", err)
				pingChannel <- Ping{ip, false}
				return
			}

			defer res.Body.Close()
			res.Close = true

			pingChannel <- Ping{ip, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-pingChannel:
			if msg.isOnline {
				fmt.Printf("Node %s is online\n", msg.ip)
			} else {
				fmt.Printf("Node %s is offline\n", msg.ip)
			}
		}
	}
}

func DoCheck(slackApi *slack.Client, channel string, ips []string) {
	checkChannel := make(chan Check)
	results := make(map[string]blockchain.CheckResponse)
	oddNodesList := OddNodes{[]string{}, []string{}, make(map[int32][]string), make(map[int32][]string), make(map[float64][]string), make(map[string][]string), make(map[string][]string)}

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s:8080/check", ip)
			response, err := http.Get(host)

			if err != nil {
				checkChannel <- Check{ip, nil, false}
				return
			}

			defer response.Body.Close()
			response.Close = true

			var check blockchain.CheckResponse

			err = json.NewDecoder(response.Body).Decode(&check)

			if err != nil {
				fmt.Printf("DoCheck: Error unmarshalling response: %v \n", err)
			}

			checkChannel <- Check{ip, &check, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-checkChannel:
			fmt.Printf("Node %s is checked\n", msg.ip)
			if msg.response != nil {
				fmt.Printf("Response: %v \n", *msg.response)
				results[msg.ip] = *msg.response
			}
		}
	}

	headers := make(map[int32][]string)
	blocks := make(map[int32][]string)
	difficulties := make(map[float64][]string)
	chainworks := make(map[string][]string)
	bestblockhashes := make(map[string][]string)

	for ip, status := range results {
		fmt.Printf("Checking IP: %s \n", ip)

		if !blockchain.IsResponseValid(status) {
			oddNodesList.RespondedWithError = append(oddNodesList.RespondedWithError, ip)
			continue
		}

		if !blockchain.DoesHeadersAndBlocksMatch(status) {
			oddNodesList.HeadersAndBlocksDontMatch = append(oddNodesList.HeadersAndBlocksDontMatch, ip)
		}

		headers[status.Status.Headers] = append(headers[status.Status.Headers], ip)
		blocks[status.Status.Blocks] = append(blocks[status.Status.Blocks], ip)
		difficulties[status.Status.Difficulty] = append(difficulties[status.Status.Difficulty], ip)
		chainworks[status.Status.Chainwork] = append(chainworks[status.Status.Chainwork], ip)
		bestblockhashes[status.Status.BestBlockHash] = append(bestblockhashes[status.Status.BestBlockHash], ip)
	}

	if len(headers) > 1 {
		maxLen := 0
		for _, ips := range headers {
			currentLength := len(ips)
			if maxLen < currentLength {
				maxLen = currentLength
			}
		}

		for key, ips := range headers {
			if len(ips) != maxLen {
				oddNodesList.HaveDifferentNumberOfHeaders[key] = ips
			}
		}
	}

	if len(blocks) > 1 {
		maxLen := 0
		for _, ips := range blocks {
			currentLength := len(ips)
			if maxLen < currentLength {
				maxLen = currentLength
			}
		}

		for key, ips := range blocks {
			if len(ips) != maxLen {
				oddNodesList.HaveDifferentNumberOfBlocks[key] = ips
			}
		}
	}

	if len(difficulties) > 1 {
		maxLen := 0
		for _, ips := range difficulties {
			currentLength := len(ips)
			if maxLen < currentLength {
				maxLen = currentLength
			}
		}

		for key, ips := range difficulties {
			if len(ips) != maxLen {
				oddNodesList.HaveDifferentDifficulty[key] = ips
			}
		}
	}

	if len(chainworks) > 1 {
		maxLen := 0
		for _, ips := range chainworks {
			currentLength := len(ips)
			if maxLen < currentLength {
				maxLen = currentLength
			}
		}

		for key, ips := range chainworks {
			if len(ips) != maxLen {
				oddNodesList.HaveDifferentChainwork[key] = ips
			}
		}
	}

	if len(bestblockhashes) > 1 {
		maxLen := 0
		for _, ips := range bestblockhashes {
			currentLength := len(ips)
			if maxLen < currentLength {
				maxLen = currentLength
			}
		}

		for key, ips := range bestblockhashes {
			if len(ips) != maxLen {
				oddNodesList.HaveDifferentBestBlockHash[key] = ips
			}
		}
	}

	PostOddNodesReport(slackApi, channel, oddNodesList)
}

func DoVersionCheck(slackApi *slack.Client, channel string, ips []string) {
	checkChannel := make(chan VersionCheck)
	results := make(map[string]blockchain.VersionResponse)
	oddNodesList := MismatchVeresionNodes{
		RespondedWithError:            []string{},
		HaveDifferentVersions:         make(map[int64][]string),
		HaveDifferentProtocolVersions: make(map[int64][]string),
		NewestVersion:                 0,
		NewestProtocolVersion:         0,
	}

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s:8080/version", ip)
			response, err := http.Get(host)

			if err != nil {
				checkChannel <- VersionCheck{ip, nil, false}
				return
			}

			defer response.Body.Close()
			response.Close = true

			var check blockchain.VersionResponse

			err = json.NewDecoder(response.Body).Decode(&check)

			if err != nil {
				fmt.Printf("DoVersionCheck: Error unmarshalling response: %v \n", err)
			}

			checkChannel <- VersionCheck{ip, &check, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-checkChannel:
			fmt.Printf("Node %s version check is receeived\n", msg.ip)
			if msg.response != nil {
				fmt.Printf("Response: %v \n", *msg.response)
				results[msg.ip] = *msg.response
			}
		}
	}

	versions := make(map[int64][]string)
	protocolVersions := make(map[int64][]string)

	for ip, status := range results {
		fmt.Printf("Checking IP: %s \n", ip)

		if !blockchain.IsVersionResponseValid(status) {
			oddNodesList.RespondedWithError = append(oddNodesList.RespondedWithError, ip)
			continue
		}

		versions[status.Info.Version] = append(versions[status.Info.Version], ip)
		protocolVersions[status.Info.ProtocolVersion] = append(protocolVersions[status.Info.ProtocolVersion], ip)
	}

	if len(versions) > 1 {
		var maxVer int64
		maxVer = 0
		for v := range versions {
			if maxVer < v {
				maxVer = v
			}
		}

		for key, ips := range versions {
			if key != maxVer {
				oddNodesList.HaveDifferentVersions[key] = ips
			}
		}

		oddNodesList.NewestVersion = maxVer
	}

	if len(protocolVersions) > 1 {
		var maxVer int64
		maxVer = 0
		for v := range protocolVersions {
			if maxVer < v {
				maxVer = v
			}
		}

		for key, ips := range protocolVersions {
			if key != maxVer {
				oddNodesList.HaveDifferentProtocolVersions[key] = ips
			}
		}

		oddNodesList.NewestProtocolVersion = maxVer
	}

	PostVersionMismatchReport(slackApi, channel, oddNodesList)
}
