package server

import "net/http"

const (
	listenAddr = ":8443"
	certPath   = "/tls/ca.pem"
	keyPath    = "/tls/key.pem"
)

type Server struct {
	server *http.Server
}

func New() *Server {
	handler := setupHandler()
	return &Server{
		server: &http.Server{
			Addr:    listenAddr,
			Handler: handler,
		},
	}
}

func (s *Server) ListenAddr(addr string) error {
	return s.server.ListenAndServeTLS(certPath, keyPath)
}
