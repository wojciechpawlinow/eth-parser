package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/wojciechpawlinow/eth-parser/internal/domain/eth"
)

const (
	url            = "https://cloudflare-eth.com"
	JSONRPCVersion = "2.0"

	// Transfer(address,address,uint256) signature
	transferHash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	// methods
	callBlockNumber       = "eth_blockNumber"
	callChainID           = "eth_chainId"
	callGetLogs           = "eth_getLogs"
	callTransactionByHash = "eth_getTransactionByHash"
)

// general structure of any req https://ethereum.org/en/developers/docs/apis/json-rpc
type ethRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int64  `json:"id"`
}

type ethResponse struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *ethError       `json:"error"`
}

type ethError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ethLog struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type ethLogResponse struct {
	Logs []ethLog `json:"result"`
}

type ethTransactionReader struct {
	block  int
	mu     sync.RWMutex
	client *http.Client
}

var _ eth.TransactionsReader = (*ethTransactionReader)(nil)

func NewEthereumTransactionsReader(client *http.Client) *ethTransactionReader {
	return &ethTransactionReader{block: 0, client: client}
}

func (r *ethTransactionReader) GetTransactions(ctx context.Context, address string) []eth.Transaction {
	r.mu.Lock()
	defer r.mu.Unlock()

	chainID, err := r.getChainID(ctx) // tip: most likely could be cached
	if err != nil {
		log.Println(fmt.Errorf("getting chain ID: %w", err))
		return []eth.Transaction{} // I'd return err, but interface forces list only
	}

	prepAddr := formatAddressForTopics(address)

	params := []any{
		map[string]any{
			"fromBlock": "0x0",
			"toBlock":   "latest",
			"address":   nil,
			"topics": []any{
				transferHash, []string{
					prepAddr,
					prepAddr,
				}},
		},
	}

	getLogsReq := &ethRequest{JSONRPC: JSONRPCVersion, Method: callGetLogs, Params: params, ID: chainID}

	resp, err := r.sendPostRequest(ctx, getLogsReq)
	if err != nil {
		log.Println(fmt.Errorf("failed to get logs: %w", err))
		return []eth.Transaction{}
	}

	var er ethLogResponse
	if err = json.Unmarshal(resp, &er); err != nil {
		log.Println(fmt.Errorf("failed json unmarshall: %w", err))
		return []eth.Transaction{}
	}

	var transactions []eth.Transaction
	for _, l := range er.Logs {

		// could go concurrent in a sync.WaitGroup
		tx, err := r.getTransactionByHash(ctx, chainID, l.TransactionHash)
		if err != nil {
			log.Println(fmt.Errorf("failed to get transaction details for hash %s: %w", l.TransactionHash, err))
			continue
		}

		transactions = append(transactions, *tx)
	}

	return transactions
}

func (r *ethTransactionReader) GetCurrentBlock(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chainID, err := r.getChainID(ctx) // tip: most likely could be cached
	if err != nil {
		log.Println(fmt.Errorf("getting chain ID: %w", err))
		return 0 // I'd return err, but interface forces int only
	}

	blockNumberReq := &ethRequest{JSONRPC: JSONRPCVersion, Method: callBlockNumber, Params: []any{}, ID: chainID}
	resp, err := r.sendPostRequest(ctx, blockNumberReq)
	if err != nil {
		log.Println(fmt.Errorf("getting block number: %w", err))
		return 0
	}

	var er ethResponse
	if err = json.Unmarshal(resp, &er); err != nil {
		log.Println(fmt.Errorf("failed json unmarshall: %w", err))
		return 0
	}

	if er.Error != nil {
		log.Println(fmt.Errorf("eth error: [%d] %s", er.Error.Code, er.Error.Message))
		return 0
	}

	var blockNumberHex string
	if err = json.Unmarshal(er.Result, &blockNumberHex); err != nil {
		log.Println(fmt.Errorf("failed to unmarshal block number: %w", err))
		return 0
	}

	blockNumber, err := hexToInt(blockNumberHex)
	if err != nil {
		log.Println(fmt.Errorf("failed converting hex to int: %w", err))
		return 0
	}

	r.block = int(blockNumber)

	return r.block
}

func (r *ethTransactionReader) getChainID(ctx context.Context) (int64, error) {
	chainIDReq := &ethRequest{JSONRPC: JSONRPCVersion, Method: callChainID, Params: []any{}, ID: 1}
	resp, err := r.sendPostRequest(ctx, chainIDReq)
	if err != nil {
		return 0, fmt.Errorf("failed sending eth request: %w", err)
	}

	var er ethResponse
	if err = json.Unmarshal(resp, &er); err != nil {
		return 0, fmt.Errorf("failed json unmarshall: %w", err)
	}

	if er.Error != nil {
		return 0, fmt.Errorf("eth error: [%d] %s", er.Error.Code, er.Error.Message)
	}

	if resp == nil || len(er.Result) < 1 {
		return 0, fmt.Errorf("missing response data")
	}

	var chainIDHex string
	if err = json.Unmarshal(er.Result, &chainIDHex); err != nil {
		return 0, fmt.Errorf("failed json unmarshall: %w", err)
	}

	id, err := hexToInt(chainIDHex)
	if err != nil {
		return 0, fmt.Errorf("failed converting hex to int: %w", err)
	}

	return id, nil
}

func (r *ethTransactionReader) getTransactionByHash(ctx context.Context, chainID int64, hash string) (*eth.Transaction, error) {
	transactionByHashReq := &ethRequest{JSONRPC: JSONRPCVersion, Method: callTransactionByHash, Params: []any{hash}, ID: chainID}

	resp, err := r.sendPostRequest(ctx, transactionByHashReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}

	var er ethResponse
	if err = json.Unmarshal(resp, &er); err != nil {
		return nil, fmt.Errorf("failed json unmarshall: %w", err)
	}

	if er.Error != nil {
		return nil, fmt.Errorf("eth error: [%d] %s", er.Error.Code, er.Error.Message)
	}

	var txDetails eth.Transaction
	if err = json.Unmarshal(er.Result, &txDetails); err != nil {
		return nil, fmt.Errorf("failed json unmarshall: %w", err)
	}

	return &txDetails, nil
}

func (r *ethTransactionReader) sendPostRequest(ctx context.Context, req *ethRequest) ([]byte, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed json marshall: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed sending http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	return respBody, nil
}

func hexToInt(hexStr string) (int64, error) {
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	return strconv.ParseInt(hexStr, 16, 64)
}

func formatAddressForTopics(address string) string {
	address = strings.ToLower(strings.TrimPrefix(address, "0x"))
	return "0x" + fmt.Sprintf("%064s", address)
}
