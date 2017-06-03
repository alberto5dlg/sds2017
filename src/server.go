package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var gUsuarios = map[string]usuario{}
var gtoken = map[string]string{}

func chkError(err error) {
	if err != nil {
		panic(err)
	}
}

func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}
	rJSON, err := json.Marshal(&r)
	chkError(err)
	w.Write(rJSON)
}

func responseJSON(w io.Writer, ok bool, info map[string]datos) {
	r := respJSON{Ok: ok, Info: info}
	rJSON, err := json.Marshal(&r)
	chkError(err)
	w.Write(rJSON)
}

func cargarBD() bool {
	decryptFile()
	raw, err := ioutil.ReadFile("bbdd.json")
	chkError(err)
	json.Unmarshal(raw, &gUsuarios)
	encryptFile()
	return true
}

func guardarBD() {
	jsonString, err := json.Marshal(gUsuarios)
	chkError(err)
	ioutil.WriteFile("bbdd.json", jsonString, 0644)
	//encryptFile()
}

//añadimos una nueva cuenta con su usuario y contraseña EJ. Facebook "username" "password"
func nuevaCuenta(usuario string, servicio string, username string, password string) {
	var newinfo datos
	newinfo.User = username
	newinfo.Pass = password
	gUsuarios[usuario].Info[servicio] = newinfo
}

func crearCuenta(resp string) bool {
	var cuen cuenta
	datos := decode64(resp)
	json.Unmarshal(datos, &cuen)
	nuevaCuenta(cuen.Boss, cuen.Servicio, cuen.User, cuen.Password)
	return true
}

func eliminarCuenta(resp string) bool {
	var cuen cuenta
	datos := decode64(resp)
	json.Unmarshal(datos, &cuen)
	delete(gUsuarios[cuen.Boss].Info, cuen.Servicio)
	return true
}

func consultarCuentas(resp string) map[string]datos {
	var cuen cuenta
	datos := decode64(resp)
	json.Unmarshal(datos, &cuen)
	return gUsuarios[cuen.Boss].Info
}

func nuevoUsuario(username string, password string, email string) {
	var newUser usuario
	hasher := sha512.Sum512([]byte(password))
	newUser.Password = encode64(hasher[:])

	newUser.Email = email
	newUser.Info = make(map[string]datos)
	newUser.Tarjetas = make(map[string]tarjeta)
	newUser.Notas = make(map[string]notas)
	gUsuarios[username] = newUser
}
func crearTarjeta(resp string) bool {
	var tar nTarjeta
	datos := decode64(resp)
	json.Unmarshal(datos, &tar)
	anyadirTarjeta(tar.Username, tar.Entidad, tar.NTarjeta, tar.Fecha, tar.CodSeg)
	return true
}

func anyadirTarjeta(username string, entidad string, nTarj string, fecha string, codSeg string) {
	var card tarjeta
	card.Entidad = entidad
	card.NTarjeta = nTarj
	card.Fecha = fecha
	card.CodSeg = codSeg
	gUsuarios[username].Tarjetas[entidad] = card
}

func crearNota(resp string) bool {
	var not nNotas
	datos := decode64(resp)
	json.Unmarshal(datos, &not)
	anyadirNotas(not.Username, not.Titulo, not.Cuerpo)
	return true
}

func anyadirNotas(username string, titulo string, cuerpo string) {
	var notes notas
	notes.Titulo = titulo
	notes.Cuerpo = cuerpo
	gUsuarios[username].Notas[titulo] = notes
}

func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	chkError(err)
	return b
}

func compLogin(resp string) string {
	var log logueado
	datos := decode64(resp)
	json.Unmarshal(datos, &log)
	hasher := sha512.Sum512([]byte(log.Password))
	password := encode64(hasher[:])

	if gUsuarios[log.User].Password == password {
		token := generarToken()
		gtoken[log.User] = token
		return token
	}
	return ""
}

func crearUsuario(resp string) string {
	var regis registrarse
	datos := decode64(resp)
	json.Unmarshal(datos, &regis)
	nuevoUsuario(regis.User, regis.Password, regis.Email)

	token := generarToken()
	gtoken[regis.User] = token
	return token
}

func generarToken() string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 16)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	switch r.Form.Get("cmd") {
	case "login":
		token := compLogin(r.Form.Get("mensaje"))
		if token == "" {
			response(w, false, "Login Erroneo")
		} else {
			response(w, true, token)
		}

	case "registro":
		token := crearUsuario(r.Form.Get("mensaje"))
		if token == "" {
			response(w, false, "Registro Erroneo")
		} else {
			response(w, true, token)
			guardarBD()
		}

	case "añadirCuenta":
		if gtoken[r.Form.Get("username")] == r.Form.Get("token") {
			if crearCuenta(r.Form.Get("mensaje")) {
				response(w, true, "Cuenta Creada")
				guardarBD()
			} else {
				response(w, false, "No se ha añadido Error")
			}
		}

	case "eliminarCuenta":
		if gtoken[r.Form.Get("username")] == r.Form.Get("token") {
			if eliminarCuenta(r.Form.Get("mensaje")) {
				response(w, true, "Cuenta Eliminada")
				guardarBD()
			} else {
				response(w, false, "No se ha eliminado Error")
			}
		}

	case "consultarCuentas":
		if gtoken[r.Form.Get("username")] == r.Form.Get("token") {
			responseJSON(w, true, consultarCuentas(r.Form.Get("mensaje")))
		}

	case "añadirTarjeta":
		if gtoken[r.Form.Get("username")] == r.Form.Get("token") {
			if crearTarjeta(r.Form.Get("mensaje")) {
				response(w, true, "Añadida la Tarjeta")
				guardarBD()
			} else {
				response(w, false, "No se ha podido añadir")
			}
		}

	case "añadirNota":
		if gtoken[r.Form.Get("username")] == r.Form.Get("token") {
			if crearNota(r.Form.Get("mensaje")) {
				response(w, true, "Añadida la Nota")
				guardarBD()
			} else {
				response(w, false, "No se ha podido añadir")
			}
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
	encryptFile()
}
