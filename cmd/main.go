package main

import (
	"go-mev/internal/strategy"
	"log"
	"math/big"
)

func main() {
	tradeSize := big.NewInt(1e18) // --> 1 ETH

	for {
		r0A, r1A := poolAState.Get()
		r0B, r1B := poolBState.Get()

		profit := strategy.SimulateArb(
			tradeSize,
			r0A, r1A,
			r0B, r1B,
		)

		if profit.Cmp(big.NewInt(0)) > 0 {
			log.Println("Arbitrage Opportunity: ", profit)
		}
	}
}
