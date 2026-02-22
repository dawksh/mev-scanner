package main

import (
	"go-mev/internal/config"
	"go-mev/internal/pool"
	"go-mev/internal/strategy"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial(config.WSURL)
	if err != nil {
		log.Fatal(err)
	}

	poolStateA := pool.NewState()
	poolStateB := pool.NewState()

	addrA := common.HexToAddress(config.PoolA)
	addrB := common.HexToAddress(config.PoolB)

	// Initialize Reserves
	err = pool.InitializeState(client, addrA, poolStateA)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.InitializeState(client, addrB, poolStateB)
	if err != nil {
		log.Fatal(err)
	}

	trigger := make(chan struct{}, 1)

	// Go routines for both pools to listen for sync event
	go pool.ListenSync(client, addrA, poolStateA, trigger)
	go pool.ListenSync(client, addrB, poolStateB, trigger)

	tradeSize := big.NewInt(1e16)

	for range trigger {
		evaluateArb(tradeSize, poolStateA, poolStateB)
	}
}

func evaluateArb(
	amountIn *big.Int,
	a, b *pool.State,
) {
	r0A, r1A := a.Get()
	r0B, r1B := b.Get()

	if r0A.Sign() == 0 || r0B.Sign() == 0 {
		return
	}
	profitAB := strategy.SimulateArb(amountIn, r0A, r1A, r0B, r1B)
	profitBA := strategy.SimulateArb(amountIn, r0B, r1B, r0A, r1A)

	if profitAB.Sign() > 0 {
		log.Println("Arb Pool A -> B:", profitAB.String())
	}

	if profitBA.Sign() > 0 {
		log.Println("Arb Pool B -> A:", profitAB.String())
	}
}
