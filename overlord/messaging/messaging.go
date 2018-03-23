package messaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meritlabs/overlord/blockchain"
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

type OddNodes struct {
	RespondedWithError           []string
	HeadersAndBlocksDontMatch    []string
	HaveDifferentNumberOfHeaders map[int32][]string
	HaveDifferentNumberOfBlocks  map[int32][]string
	HaveDifferentDifficulty      map[float64][]string
	HaveDifferentChainwork       map[string][]string
	HaveDifferentBestBlockHash   map[string][]string
}

func DoPing(ips []string) {
	pingChannel := make(chan Ping)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s:8080/ping", ip)
			_, err := http.Get(host)

			if err != nil {
				pingChannel <- Ping{ip, false}
				return
			}

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
			defer response.Body.Close()

			if err != nil {
				checkChannel <- Check{ip, nil, false}
				return
			}

			var check blockchain.CheckResponse

			err = json.NewDecoder(response.Body).Decode(&check)

			if err != nil {
				fmt.Printf("DoCheck: Error unmarshalling response: %v \n", err)
			}

			fmt.Printf("Check result is: %v \n", check)

			checkChannel <- Check{ip, &check, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-checkChannel:
			fmt.Printf("Node %s is checked\n", msg.ip)
			results[msg.ip] = *msg.response
		}
	}

	fmt.Printf("Results: %v \n", results)

	headers := make(map[int32][]string)
	blocks := make(map[int32][]string)
	difficulties := make(map[float64][]string)
	chainworks := make(map[string][]string)
	bestblockhashes := make(map[string][]string)

	for ip, status := range results {
		fmt.Printf("Checking IP: %s \n", ip)

		if !blockchain.IsResponseValid(status) {
			fmt.Printf("Errored! \n")
			oddNodesList.RespondedWithError = append(oddNodesList.RespondedWithError, ip)
			continue
		}

		if !blockchain.DoesHeadersAndBlocksMatch(status) {
			fmt.Printf("Heanders and Blocks: %b \n", blockchain.DoesHeadersAndBlocksMatch(status))
			// Write error to slack
			oddNodesList.HeadersAndBlocksDontMatch = append(oddNodesList.HeadersAndBlocksDontMatch, ip)
		}

		headers[status.Status.Headers] = append(headers[status.Status.Headers], ip)
		blocks[status.Status.Blocks] = append(blocks[status.Status.Blocks], ip)
		difficulties[status.Status.Difficulty] = append(difficulties[status.Status.Difficulty], ip)
		chainworks[status.Status.Chainwork] = append(chainworks[status.Status.Chainwork], ip)
		bestblockhashes[status.Status.BestBlockHash] = append(bestblockhashes[status.Status.BestBlockHash], ip)
	}

	fmt.Printf("Headers: %v \n", headers)

	if len(headers) > 1 {
		fmt.Printf("Cluster have differrent number of headers on different nodes")
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

	fmt.Printf("Blocks: %v \n", blocks)

	if len(blocks) > 1 {
		fmt.Printf("Cluster have differrent number of blocks on different nodes")

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

	fmt.Printf("Difficulty: %v \n", difficulties)

	if len(difficulties) > 1 {
		fmt.Printf("Cluster have differrent difficulty on different nodes")

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

	fmt.Printf("Chainwork: %v \n", chainworks)

	if len(chainworks) > 1 {
		fmt.Printf("Cluster have differrent chainwork state on different nodes")

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

	fmt.Printf("BestBlockHash: %v \n", bestblockhashes)

	if len(bestblockhashes) > 1 {
		fmt.Printf("Cluster have differrent BestBlockHash on different nodes")
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
