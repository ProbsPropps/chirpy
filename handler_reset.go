package main

import "net/http"


func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Access denied", nil)
	}
	cfg.queries.DeleteUsers(req.Context())
}
