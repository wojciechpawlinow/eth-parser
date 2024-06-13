package memory

import (
	"context"
	"sync"

	"github.com/wojciechpawlinow/eth-parser/internal/domain/parser"
)

type addressRepository struct {
	addresses map[string]struct{}
	mu        sync.RWMutex
}

var _ parser.AddressRepository = (*addressRepository)(nil)

func NewAddressRepository() *addressRepository {
	return &addressRepository{addresses: make(map[string]struct{})}
}

func (r *addressRepository) Add(_ context.Context, address string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.addresses[address] = struct{}{}
}

func (r *addressRepository) IsSubscribed(_ context.Context, address string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.addresses[address]

	return ok
}
