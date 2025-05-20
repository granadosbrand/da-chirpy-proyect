package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, 403, "Wrong platform:", errors.New(cfg.platform))
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error deleting all users", err)
		return
	}

	respondWithJSon(w, 200, struct {
		OK bool `json:"ok"`
	}{
		OK: true,
	})
	// cfg.FileServerHits.Store(0)
	// w.Write([]byte("Counter reset successfully"))
}
