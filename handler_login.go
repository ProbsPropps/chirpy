package main

import (
	"encoding/json"
	"net/http"

	"github.com/ProbsPropps/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode", err)
	}

	user, err := cfg.queries.GetUserByEmail(req.Context(), params.Email)
	
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
