package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var gUsuarios = map[string]usuario{}

type login struct {
	User     string
	Password string
}

type registro struct {
	User     string
	Password string
	Email    string
}

type datos struct {
	User string
	Pass string
}

type usuario struct {
	Email    string
	Password string
	Info     map[string]datos
}

type resp struct {
	Ok  bool
	Msg string
}

func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}
	rJSON, err := json.Marshal(&r)
	if err != nil {
		panic(err)
	}
	w.Write(rJSON)
}

func cargarBD() bool {
	raw, err := ioutil.ReadFile("bbdd.json")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	json.Unmarshal(raw, &gUsuarios)
	return true
}

func guardarBD() {
	jsonString, err := json.Marshal(gUsuarios)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("bbdd.json", jsonString, 0644)
}

//añadimos una nueva cuenta con su usuario y contraseña EJ. Facebook "username" "password"
func nuevaCuenta(usuario string, servicio string, username string, password string) {
	var newinfo datos
	newinfo.User = username
	newinfo.Pass = password
	gUsuarios[usuario].Info[servicio] = newinfo
}

func nuevoUsuario(username string, password string, email string) {
	var newUser usuario
	newUser.Password = password
	newUser.Email = email
	newUser.Info = make(map[string]datos)
	gUsuarios[username] = newUser

}

func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func compLogin(resp string) bool {
	var log login
	datos := decode64(resp)
	json.Unmarshal(datos, &log)
	if gUsuarios[log.User].Password == log.Password {
		return true
	}
	return false
}

func crearUsuario(resp string) bool {
	var regis registro
	datos := decode64(resp)
	json.Unmarshal(datos, &regis)
	nuevoUsuario(regis.User, regis.Password, regis.Email)
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	switch r.Form.Get("cmd") {
	case "login":
		if compLogin(r.Form.Get("mensaje")) {
			response(w, true, "Login Correcto")
		} else {
			response(w, false, "Login Erroneo")
		}

	case "registro":
		if crearUsuario(r.Form.Get("mensaje")) {
			response(w, true, "Usuario Creado")
		}

	}
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
	cargarBD()
	conectServer()
	guardarBD()
}
