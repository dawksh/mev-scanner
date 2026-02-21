package strategy

import (
	"go-mev/internal/math"
	"math/big"
)

func SimulateArb(
	amountIn *big.Int,
	r0A, r1A *big.Int,
	r0B, r1B *big.Int,
) *big.Int {
	// Swap on A
	outA := math.GetAmountOut(amountIn, r0A, r1A)

	// Swap on B
	outB := math.GetAmountOut(outA, r1B, r0B)

	return new(big.Int).Sub(outB, amountIn)
}
