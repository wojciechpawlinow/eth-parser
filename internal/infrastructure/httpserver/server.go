package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/wojciechpawlinow/eth-parser/internal/infrastructure/container"
)

type Server struct {
	*http.Server
	ctn *container.Container
}

// Run is a Server constructor that starts the HTTP server in a goroutine and enables routing
func Run(ctn *container.Container, errChan chan error) *Server {

	addr := fmt.Sprintf("0.0.0.0:8080")

	s := &Server{
		&http.Server{
			Addr: addr, // config
		},
		ctn,
	}

	http.HandleFunc("/current-block", s.handleGetCurrentBlock)
	http.HandleFunc("/subscribe", s.handleSubscribe)
	http.HandleFunc("/address/", s.handleGetTransactions)

	go func() {
		log.Println(fmt.Sprintf("listening at %s", addr))
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return s
}

// Shutdown is a Shutdown function overload
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.Server.Shutdown(ctx)
}
