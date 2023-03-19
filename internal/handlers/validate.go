package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/thevilledev/learn-admission-controllers/pkg/admission"
)

func ValidateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := validateFunc(w, r); err != nil {
			log.Printf("Could not handle the request: %+v", err)
		}
	})
}

func validateFunc(w http.ResponseWriter, r *http.Request) error {
	log.Print("Handling webhook request ...")

	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid method %s, only POST requests are allowed", r.Method)
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "invalid content-type", http.StatusBadRequest)
		return fmt.Errorf("unsupported content type")
	}

	defer r.Body.Close()
	lr := &io.LimitedReader{
		R: r.Body,
		N: int64(3*1024*1024) + 1, // 3 MB
	}
	body, err := io.ReadAll(lr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("could not read request body: %v", err)
	}
	if lr.N <= 0 {
		http.Error(w, "too large", http.StatusRequestHeaderFieldsTooLarge)
		return fmt.Errorf("entity too large")
	}

	bytes, err := admission.Admit(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("failed to admit: %s", err.Error())
	}

	log.Print("Webhook request handled successfully")
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Could not write response: %+v", err)
	}

	return nil
}
