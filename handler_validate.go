package main

import (
	"net/http"
	"encoding/json"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {                                    
		Body string `json:"body"`                               
	}                                                           
        
	type returnValid struct {
		Valid bool `json:"valid"`
	}
                                                                
	decoder := json.NewDecoder(req.Body)                        
	params := parameters{}                                      
	err := decoder.Decode(&params)                              
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	}

	respondWithJSON(w, http.StatusOK, returnValid{
		Valid: true,
	})
}
