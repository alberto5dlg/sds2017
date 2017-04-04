package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type datos struct {
	User string
	Pass string
}

type usuario struct {
	Email    string
	Password string
	Info     map[string]datos
}

func cargarBD(gUsuarios map[string]usuario) bool {
	raw, err := ioutil.ReadFile("bbdd.json")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	json.Unmarshal(raw, &gUsuarios)
	return true
}

func guardarBD(gUsuarios map[string]usuario) {
	jsonString, err := json.Marshal(gUsuarios)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("bbdd.json", jsonString, 0644)
}

//añadimos una nueva cuenta con su usuario y contraseña EJ. Facebook "username" "password"
func nuevaCuenta(usuario string, servicio string, username string, password string, gUsuarios map[string]usuario) {
	var newinfo datos
	newinfo.User = username
	newinfo.Pass = password
	gUsuarios[usuario].Info[servicio] = newinfo
}

func nuevoUsuario() {

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}
func conectServer() {
	stopChan := make(chan os.Signal)
	log.Println("Escuchando en: 127.0.0.1:8081 ... ")
	signal.Notify(stopChan, os.Interrupt)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler))
	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", mux)
	<-stopChan
	log.Println("Apagando servidor ...")
	log.Println("Servidor detenido correctamente")
}

func main() {
	gUsuarios := make(map[string]usuario)
	cargarBD(gUsuarios)
	nuevaCuenta("Fer", "Instagram", "ferchu", "adios", gUsuarios)
	guardarBD(gUsuarios)
	//fmt.Println(gUsuarios["Fer"].Info["Instagram"])
	conectServer()

}
