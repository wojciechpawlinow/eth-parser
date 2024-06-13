package parser

import "context"

type AddressRepository interface {
	// Add subscribe new blockchain address
	Add(ctx context.Context, address string)
	// IsSubscribed check if is given address is observed
	IsSubscribed(ctx context.Context, address string) bool
}
