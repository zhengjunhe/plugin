package deploy

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
)


func setupEthClient(ethURL string) (*ethclient.Client, error) {
	if strings.TrimSpace(ethURL) == "" {
		return nil, nil
	}
	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket client %w", err)
	}

	return client, nil
}

func main() {
	setupEthClient
}
