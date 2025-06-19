package main

import "net/http"

func main() {
	mu := http.NewServeMux()
	mu.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Handler: mu,
		Addr: ":8080",
	}
	server.ListenAndServe()
}
