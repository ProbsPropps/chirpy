package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ProbsPropps/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	User_ID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {                                    
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}                                                           
        
	type returnValid struct {
		CleanedBody string `json:"cleaned_body"`
	}
                                                                
	decoder := json.NewDecoder(req.Body)                        
	params := parameters{}                                      
	err := decoder.Decode(&params)                              
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	err = validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	cleanChirp := cleanBody(params.Body)
	chirp, err := cfg.queries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanChirp,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Trouble with CreateChirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_ID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.queries.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
	}
	
	timeline := make([]Chirp, len(chirps))

	for i, chirp := range chirps {
		data := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_ID: chirp.UserID,
		}
		timeline[i] = data
	}

	respondWithJSON(w, http.StatusOK, timeline)
}
