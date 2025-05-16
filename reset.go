package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits.Store(0)
	w.Write([]byte("Counter reset successfully"))
}