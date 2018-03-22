package blockchain

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

var (
	binary    = "/usr/local/bin/merit-cli"
	mode      = "mainnet"
	statusCmd = "getblockchaininfo"
)

/*
{
  "chain": "main",
  "blocks": 117092,
  "headers": 117092,
  "bestblockhash": "07eb1c7e6213a112635a329a7fccedf389c646438e1fe665220c7787587bb91e",
  "difficulty": 1.918193663961827e-08,
  "mediantime": 1521716735,
  "verificationprogress": 1,
  "chainwork": "00000000000000000000000000000000000000000000000000000000004173ba",
  "pruned": false
}
*/

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int32   `json:"blocks"`
	Headers              int32   `json:"headers"`
	BestBlockHash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Mediantime           int32   `json:"mediantime"`
	VerificationProgress int     `json:"verificationprogress"`
	Chainwork            string  `json:"chainwork"`
	Pruned               bool    `json:"pruned"`
}

type CheckResponse struct {
	Status BlockchainInfo `json:"status"`
	Error  error          `json:"error"`
}

func GetBlockchainInfo() (*BlockchainInfo, error) {
	res, err := execute(binary, "--"+mode, statusCmd)

	if err != nil {
		fmt.Printf("Error executing getblockchaininfo command: %v \n", err)
		return nil, err
	}

	var info BlockchainInfo

	err = json.Unmarshal(res, &info)

	if err != nil {
		fmt.Printf("Error unmarshalling response: %v \n", err)
		return nil, err
	}

	return &info, nil
}

func execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	return cmd.Output()
}

func IsResponseValid(res CheckResponse) bool {
	return res.Error == nil
}

func DoesHeadersAndBlocksMatch(res CheckResponse) bool {
	return res.Status.Headers == res.Status.Blocks
}
