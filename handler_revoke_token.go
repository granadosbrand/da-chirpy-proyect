package main

import (
	"net/http"

	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error extracting token", err)
		return
	}

	// Revoke token
	_, err = cfg.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}

	respondWithJSon(w, http.StatusNoContent, nil)
}
