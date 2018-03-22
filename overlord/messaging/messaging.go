package messaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meritlabs/overlord/blockchain"
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

var (
	ips = []string{"127.0.0.1:8080", "192.168.1.6:8080"} // grab it from config
)

func DoPing() {
	pingChannel := make(chan Ping)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s/ping", ip)
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

func DoCheck() {
	checkChannel := make(chan Check)
	results := make(map[string]blockchain.CheckResponse)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s/check", ip)
			response, err := http.Get(host)
			defer response.Body.Close()

			if err != nil {
				checkChannel <- Check{ip, nil, false}
				return
			}

			var check blockchain.CheckResponse

			err = json.NewDecoder(response.Body).Decode(&check)

			if err != nil {
				fmt.Printf("Error unmarshalling response: %v \n", err)
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
			// Write error to slack
			continue
		}

		if !blockchain.DoesHeadersAndBlocksMatch(status) {
			fmt.Printf("Heanders and Blocks: %b \n", blockchain.DoesHeadersAndBlocksMatch(status))
			// Write error to slack
		}

		headers[status.Status.Headers] = append(headers[status.Status.Headers], ip)
		blocks[status.Status.Blocks] = append(blocks[status.Status.Blocks], ip)
		difficulties[status.Status.Difficulty] = append(difficulties[status.Status.Difficulty], ip)
		chainworks[status.Status.Chainwork] = append(chainworks[status.Status.Chainwork], ip)
		bestblockhashes[status.Status.BestBlockHash] = append(bestblockhashes[status.Status.BestBlockHash], ip)
	}

	fmt.Printf("Headers: %v \n", headers)
	fmt.Printf("Blocks: %v \n", blocks)
	fmt.Printf("Difficulty: %v \n", difficulties)
	fmt.Printf("Chainwork: %v \n", chainworks)
	fmt.Printf("BestBlockHash: %v \n", bestblockhashes)
}
