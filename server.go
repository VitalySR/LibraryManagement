package LibraryManagement

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler, readTimeoutSeconds int, writeTimeoutSeconds int) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    time.Duration(readTimeoutSeconds) * time.Second,
			WriteTimeout:   time.Duration(writeTimeoutSeconds) * time.Second,
		},
	}
}

func (s *Server) Run() error {
	log.Println("Server is running on " + s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server is shutting down " + s.httpServer.Addr)
	return s.httpServer.Shutdown(ctx)
}
