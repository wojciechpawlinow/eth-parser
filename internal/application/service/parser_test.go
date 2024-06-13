package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/eth-parser/internal/domain/eth"
)

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) IsSubscribed(ctx context.Context, address string) bool {
	args := m.Called(ctx, address)
	return args.Bool(0)
}

func (m *MockAddressRepository) Add(ctx context.Context, address string) {
	m.Called(ctx, address)
}

type MockTransactionsReader struct {
	mock.Mock
}

func (m *MockTransactionsReader) GetCurrentBlock(ctx context.Context) int {
	args := m.Called(ctx)
	return args.Int(0)
}

func (m *MockTransactionsReader) GetTransactions(ctx context.Context, address string) []eth.Transaction {
	args := m.Called(ctx, address)
	return args.Get(0).([]eth.Transaction)
}

func TestSubscribe(t *testing.T) {
	mockAddressDb := new(MockAddressRepository)
	mockEth := new(MockTransactionsReader)
	prs := NewParser(mockAddressDb, mockEth)

	ctx := context.Background()
	address := "0x123"

	mockAddressDb.On("IsSubscribed", ctx, address).Return(false)
	mockAddressDb.On("Add", ctx, address).Return()

	result := prs.Subscribe(address)

	assert.True(t, result)
	mockAddressDb.AssertCalled(t, "Add", ctx, address)
}

func TestSubscribeAlreadySubscribed(t *testing.T) {
	mockAddressDb := new(MockAddressRepository)
	mockEth := new(MockTransactionsReader)
	prs := NewParser(mockAddressDb, mockEth)

	ctx := context.Background()
	address := "0x123"

	mockAddressDb.On("IsSubscribed", ctx, address).Return(true)

	result := prs.Subscribe(address)

	assert.True(t, result)
	mockAddressDb.AssertNotCalled(t, "Add", ctx, address)
}

func TestGetCurrentBlock(t *testing.T) {
	mockAddressDb := new(MockAddressRepository)
	mockEth := new(MockTransactionsReader)
	prs := NewParser(mockAddressDb, mockEth)

	ctx := context.Background()
	blockNumber := 12345

	mockEth.On("GetCurrentBlock", ctx).Return(blockNumber)

	result := prs.GetCurrentBlock()

	assert.Equal(t, blockNumber, result)
	mockEth.AssertCalled(t, "GetCurrentBlock", ctx)
}

func TestGetTransactions(t *testing.T) {
	mockAddressDb := new(MockAddressRepository)
	mockEth := new(MockTransactionsReader)
	prs := NewParser(mockAddressDb, mockEth)

	ctx := context.Background()
	address := "0x123"
	transactions := []eth.Transaction{{Hash: "0xabc"}}

	mockAddressDb.On("IsSubscribed", ctx, address).Return(true)
	mockEth.On("GetTransactions", ctx, address).Return(transactions)

	result := prs.GetTransactions(address)

	assert.Equal(t, transactions, result)
	mockEth.AssertCalled(t, "GetTransactions", ctx, address)
}

func TestGetTransactionsNotSubscribed(t *testing.T) {
	mockAddressDb := new(MockAddressRepository)
	mockEth := new(MockTransactionsReader)
	prs := NewParser(mockAddressDb, mockEth)

	ctx := context.Background()
	address := "0x123"

	mockAddressDb.On("IsSubscribed", ctx, address).Return(false)

	result := prs.GetTransactions(address)

	assert.Empty(t, result)
	mockEth.AssertNotCalled(t, "GetTransactions", ctx, address)
}
