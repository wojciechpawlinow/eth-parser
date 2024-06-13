package eth

import "context"

type Transaction struct {
	BlockHash        *string `json:"-"`
	BlockNumber      *string `json:"-"`
	TransactionIndex *string `json:"transactionIndex"`
	Hash             string  `json:"hash"`
	From             string  `json:"from"`
	To               *string `json:"to"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Input            string  `json:"-"`
	Nonce            string  `json:"-"`
	Value            string  `json:"value"`
	V                string  `json:"-"`
	R                string  `json:"-"`
	S                string  `json:"-"`
}

type TransactionsReader interface {
	// GetTransactions fetch list of transactions for a subscribed address
	GetTransactions(ctx context.Context, address string) []Transaction
	// GetCurrentBlock last parsed block number
	GetCurrentBlock(ctx context.Context) int
}
