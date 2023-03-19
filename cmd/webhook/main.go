package main

import (
	"log"
	"os"

	"github.com/thevilledev/learn-admission-controllers/internal/server"
)

const (
	defaultListenAddr = ":8443"
	defaultCertPath   = "/run/secrets/tls/tls.crt"
	defaultKeyPath    = "/run/secrets/tls/tls.key"
)

func getenv(k, d string) string {
	v := os.Getenv(k)
	if len(v) == 0 {
		return d
	}
	return v
}

func main() {
	listenAddr := getenv("LISTEN_ADDR", defaultListenAddr)
	certPath := getenv("CERT_PATH", defaultCertPath)
	keyPath := getenv("KEY_PATH", defaultKeyPath)

	s, err := server.New(listenAddr, certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
