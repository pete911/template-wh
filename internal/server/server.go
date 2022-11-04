package server

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Mutate func([]byte) ([]byte, error)

func ListenAndServeTLS(mutateFn Mutate, certFile, keyFile string) error {

	mux := http.NewServeMux()
	mux.Handle("/mutate", &mutateHandler{mutateFn: mutateFn})

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}
	return s.ListenAndServeTLS(certFile, keyFile)
}

type mutateHandler struct {
	mutateFn Mutate
}

func (m *mutateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("cannot read body: %v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	mutated, err := m.mutateFn(body)
	if err != nil {
		log.Printf("cannot mutate request: %v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(mutated); err != nil {
		log.Printf("cannot write response: %v", err)
	}
}
