package models

import "math/big"

type MintInfo struct {
	ControllerAddress 	string
	RecipientAddress 	string
	Amount 				*big.Int
	TotalMint 			*big.Int
}