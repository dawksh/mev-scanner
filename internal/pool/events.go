package pool

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ListenSync(
	client *ethclient.Client,
	poolAddr common.Address,
	state *State,
) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{poolAddr},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(
		context.Background(),
		query,
		logs,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)

			case vLog := <-logs:
				// Decode sync event
			}
		}
	}()
}
