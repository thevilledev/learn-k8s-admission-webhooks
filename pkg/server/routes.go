package server

import (
	"net/http"

	"github.com/thevilledev/learn-admission-controllers/pkg/handlers"
)

func setupHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/mutate", handlers.MutateHandler())
	//mux.Handle("/validate", webhookHandler(validateFunc))
	return mux
}

/*func validateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateFunc(w, r)
	})
}*/
