package math

import (
	"math/big"
)

var (
	feeMul = big.NewInt(997)
	feeDiv = big.NewInt(1000)
)

func GetAmountOut(amountIn, reserveIn, reserveOut *big.Int) *big.Int {
	amountInWithFee := new(big.Int).Mul(amountIn, feeMul)

	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)

	denominator := new(big.Int).Add(
		new(big.Int).Mul(reserveIn, feeDiv),
		amountInWithFee,
	)

	if denominator.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}

	return numerator.Div(numerator, denominator)
}
