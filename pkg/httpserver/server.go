package httpserver

import (
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(handler http.Handler, opts ...Options) *Server {
	httpServer := &http.Server{
		Handler: handler,
	}

	server := &Server{
		httpServer,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Startup() error {
	return s.ListenAndServe()
}
