package helpers

import (
	"fmt"
	"math/big"
)

func StringToBigFloat(input string) (*big.Float, error) {

	var price big.Float

	if _, ok := price.SetString(input); !ok {
		return big.NewFloat(0), fmt.Errorf("invalid number format: %s", input)
	}

	return &price, nil
}
