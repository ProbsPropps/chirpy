package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ProbsPropps/chirpy/internal/auth"
	"github.com/ProbsPropps/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode email", err)
		return
	}
	user, err := cfg.queries.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user from email", err)
		return
	}
	hashPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
	}
	
	if err = cfg.queries.AddHashPass(req.Context(), database.AddHashPassParams{
		HashedPassword: sql.NullString{String: hashPass, Valid: true},
		Email: params.Email,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add hashed password to database", err)
	}
	
	respondWithJSON(w, http.StatusCreated, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
