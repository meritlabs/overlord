package blockchain

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const (
	binary      = "/usr/local/bin/merit-cli"
	statusCmd   = "getblockchaininfo"
	veresionCmd = "getnetworkinfo"
)

// BlockchainInfo - blockchaininfo response data
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

// CheckResponse - status http response
type CheckResponse struct {
	Status BlockchainInfo `json:"status"`
	Error  string         `json:"error"`
}

// VersionInfo - daemon and protocol level data
type VersionInfo struct {
	Version         int64 `json:"version"`
	ProtocolVersion int64 `json:"protocolversion"`
}

// VersionResponse - version response
type VersionResponse struct {
	Info  VersionInfo `json:"info"`
	Error string      `json:"error"`
}

// GetBlockchainInfo - executes getblockchaininfo command
func GetBlockchainInfo(mode string) (*BlockchainInfo, error) {
	res, err := execute(binary, "--"+mode, statusCmd)

	if err != nil {
		fmt.Printf("GetBlockchainInfo: Error executing getblockchaininfo command: %v \n", err)
		return nil, err
	}

	var info BlockchainInfo

	err = json.Unmarshal(res, &info)

	if err != nil {
		fmt.Printf("GetBlockchainInfo: Error unmarshalling response: %v \n", err)
		return nil, err
	}

	return &info, nil
}

// GetInfo - executes getinfo command
func GetInfo(mode string) (*VersionInfo, error) {
	res, err := execute(binary, "--"+mode, veresionCmd)

	if err != nil {
		fmt.Printf("GetInfo: Error executing getinfo command: %v \n", err)
		return nil, err
	}

	fmt.Printf("Info %s \n", string(res[:len(res)]))

	var info VersionInfo

	err = json.Unmarshal(res, &info)
	fmt.Printf("Info %v \n", info)

	if err != nil {
		fmt.Printf("GetInfo: Error unmarshalling response: %v \n", err)
		return nil, err
	}

	return &info, nil
}

func execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	return cmd.Output()
}

// IsResponseValid - checks CheckResponse for errors
func IsResponseValid(res CheckResponse) bool {
	return res.Error == ""
}

// IsVersionResponseValid - checks VersionResponse for errors
func IsVersionResponseValid(res VersionResponse) bool {
	return res.Error == ""
}

// DoesHeadersAndBlocksMatch - checks CheckResponse for errors
func DoesHeadersAndBlocksMatch(res CheckResponse) bool {
	return res.Status.Headers == res.Status.Blocks
}
