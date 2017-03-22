package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://127.0.0.1:8081"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	log.Println("Escuchando en: 127.0.0.1:8081 ... ")
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler))

	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", mux)
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// apagar servidor de forma segura

	log.Println("Servidor detenido correctamente")

}
