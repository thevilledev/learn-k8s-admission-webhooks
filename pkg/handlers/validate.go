package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/thevilledev/learn-admission-controllers/pkg/admission/validate"
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid method %s, only POST requests are allowed", r.Method)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("could not read request body: %v", err)
	}

	if contentType := r.Header.Get("Content-Type"); contentType != jsonContentType {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("unsupported content type %s, only %s is supported", contentType, jsonContentType)
	}

	var writeErr error
	if bytes, err := validate.Func(body); err != nil {
		log.Printf("Error handling webhook request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr = w.Write([]byte(err.Error()))
	} else {
		log.Print("Webhook request handled successfully")
		_, writeErr = w.Write(bytes)
	}

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
	return nil
}
