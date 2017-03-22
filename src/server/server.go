package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://127.0.0.1:8081"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/", handler)
	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
