package pool

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const pairABIJSON = `[{
	"constant":true,
	"inputs":[],
	"name":"getReserves",
	"outputs":[
		{"internalType":"uint112","name":"reserve0","type":"uint112"},
		{"internalType":"uint112","name":"reserve1","type":"uint112"},
		{"internalType":"uint32","name":"blockTimestampLast","type":"uint32"}
	],
	"type":"function"
},{
	"anonymous":false,
	"inputs":[
		{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},
		{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}
	],
	"name":"Sync",
	"type":"event"
}]`

var syncTopic = common.HexToHash(
	"0x1c411e9a96e0b58e4f3f0a5b2b5c2e5e1e8b1b9eebf5f9d8e9e9e9e9e9e9e9e9", // Sync event topic hash
)

func InitializeState(
	client *ethclient.Client,
	addr common.Address,
	state *State,
) error {
	parsedABI, err := abi.JSON(strings.NewReader(pairABIJSON))
	if err != nil {
		return err
	}

	calldata, err := parsedABI.Pack("getReserves")
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		To:   &addr,
		Data: calldata,
	}

	res, err := client.CallContract(
		context.Background(),
		msg,
		nil,
	)
	if err != nil {
		return err
	}

	var output struct {
		Reserve0 *big.Int
		Reserve1 *big.Int
		Time     uint32
	}
	err = parsedABI.UnpackIntoInterface(
		&output,
		"getReserves",
		res,
	)
	if err != nil {
		return err
	}

	state.Update(output.Reserve0, output.Reserve1)

	return nil
}

func ListenSync(
	client *ethclient.Client,
	poolAddr common.Address,
	state *State,
	trigger chan struct{},
) {

	parsedABI, err := abi.JSON(strings.NewReader(pairABIJSON))
	if err != nil {
		log.Fatal(err)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{poolAddr},
		Topics:    [][]common.Hash{{syncTopic}},
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

	log.Println("Listening to sync events")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)

		case vLog := <-logs:

			var event struct {
				Reserve0 *big.Int
				Reserve1 *big.Int
			}

			err := parsedABI.UnpackIntoInterface(
				&event,
				"Sync",
				vLog.Data,
			)
			if err != nil {
				log.Fatal(err)
			}

			state.Update(event.Reserve0, event.Reserve1)
			select {
			case trigger <- struct{}{}:
			default:
			}
		}
	}
}
