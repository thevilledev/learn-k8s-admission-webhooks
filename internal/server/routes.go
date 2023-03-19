package server

import (
	"net/http"

	"github.com/thevilledev/learn-admission-controllers/internal/handlers"
)

func setupHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/validate", handlers.ValidateHandler())
	return mux
}
