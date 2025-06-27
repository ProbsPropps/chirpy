package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {                                    
		Body string `json:"body"`                               
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

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	
	respondWithJSON(w, http.StatusOK, returnValid{
		CleanedBody: cleanBody(params.Body),
	})
}

func cleanBody(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if slices.Contains(badWords, lowerWord){
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return strings.TrimLeft(cleaned, " ")
}
