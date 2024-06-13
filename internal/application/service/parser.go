package service

import (
	"context"
	"fmt"
	"log"

	"github.com/wojciechpawlinow/eth-parser/internal/domain/eth"
	"github.com/wojciechpawlinow/eth-parser/internal/domain/parser"
)

type parserPort struct {
	addressDb parser.AddressRepository
	eth       eth.TransactionsReader
}

var _ parser.Parser = (*parserPort)(nil)

func NewParser(addressDb parser.AddressRepository, eth eth.TransactionsReader) *parserPort {
	return &parserPort{
		addressDb: addressDb,
		eth:       eth,
	}
}

func (p *parserPort) Subscribe(address string) bool {
	ctx := context.Background() // that should come as a param

	if !p.addressDb.IsSubscribed(ctx, address) {
		p.addressDb.Add(ctx, address)
		return true
	}

	log.Println(fmt.Sprintf("address %s already has a subscription", address))
	return true
}

func (p *parserPort) GetCurrentBlock() int {
	ctx := context.Background() // that should come as a param

	return p.eth.GetCurrentBlock(ctx)
}

func (p *parserPort) GetTransactions(address string) []eth.Transaction {
	ctx := context.Background() // that should come as a param

	if p.addressDb.IsSubscribed(ctx, address) {
		return p.eth.GetTransactions(ctx, address)
	}

	log.Println(fmt.Sprintf("address %s has no subscription, return empty list", address))
	return []eth.Transaction{}
}
