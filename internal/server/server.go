package server

import (
	"crypto/tls"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(addr string, certPath string, keyPath string) (*Server, error) {
	handler := setupHandler()

	cer, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cer},
			},
		},
	}, nil
}

func (s *Server) Listen() error {
	return s.server.ListenAndServeTLS("", "")
}
