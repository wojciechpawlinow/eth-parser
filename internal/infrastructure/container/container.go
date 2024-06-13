package container

import (
	"net/http"
	"time"

	"github.com/wojciechpawlinow/eth-parser/internal/application/service"
	"github.com/wojciechpawlinow/eth-parser/internal/domain/parser"
	"github.com/wojciechpawlinow/eth-parser/internal/infrastructure/blockchain"
	"github.com/wojciechpawlinow/eth-parser/internal/infrastructure/database/memory"
)

type Container struct {
	ParserPort parser.Parser
}

func New() *Container {
	addressDb := memory.NewAddressRepository()
	ethReader := blockchain.NewEthereumTransactionsReader(&http.Client{Timeout: 3 * time.Minute})
	parserPort := service.NewParser(addressDb, ethReader)

	return &Container{
		ParserPort: parserPort,
	}
}
