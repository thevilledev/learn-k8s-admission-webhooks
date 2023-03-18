package server

import (
	"net/http"

	"github.com/thevilledev/learn-admission-controllers/pkg/handlers"
)

func setupHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/mutate", handlers.MutateHandler())
	mux.Handle("/validate", handlers.ValidateHandler())
	return mux
}
