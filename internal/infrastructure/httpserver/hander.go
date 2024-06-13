package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type subscribeRequestModel struct {
	Address string `json:"address"`
}

//	POST /subscribe
//
// { "address": "" }
func (srv *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req subscribeRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("%s: %s address: %s\n", http.MethodPost, r.URL.Path, req.Address)

	if req.Address == "" {
		http.Error(w, "address field is required", http.StatusBadRequest)
		return
	}

	if !srv.ctn.ParserPort.Subscribe(req.Address) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// GET /current-block
func (srv *Server) handleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("%s: %s\n", http.MethodGet, r.URL.Path)

	response := map[string]int{
		"current_block": srv.ctn.ParserPort.GetCurrentBlock(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GET /transactions/{address}
func (srv *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("%s: %s\n", http.MethodGet, r.URL.Path)

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[1] != "address" || parts[3] != "transactions" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	address := parts[2]

	// here we should handle a non subscribed address, that should return err
	// to simplify, if there is no subscription then list is empty
	transactions := srv.ctn.ParserPort.GetTransactions(address)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transactions)
}
