package main

import (
	"errors"
	"slices"
	"strings"
)

func validateChirp(msg string) error {
	const maxChirpLength = 140
	if len(msg) > maxChirpLength {
		return errors.New("Max 140 Characters Allowed")
	}
	return nil
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
