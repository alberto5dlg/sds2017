package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var urlServer = "https://127.0.0.1:8081"

type datos struct {
	User string
	Pass string
}

type userRes struct {
	User     string
	Password string
}

type cuentaRes struct {
	Boss     string
	Servicio string
	User     string
	Password string
}

type structUser struct {
	User     string
	Password string
	Email    string
}

type resp struct {
	Ok  bool
	Msg string
}

type respJSON struct {
	Ok   bool
	Info map[string]datos
}

func chkError(err error) {
	if err != nil {
		panic(err)
	}
}

func menu() int {
	var opcion = 0
	for opcion <= 0 || opcion >= 4 {
		fmt.Printf("Aplicación SDS Seguridad\n")
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("1 - Login\n")
		fmt.Printf("2 - Registro\n")
		fmt.Printf("3 - Salir\n")
		fmt.Printf("Opción: ")
		fmt.Scanf("%d\n", &opcion)
	}
	return opcion
}

func ignorarHTTPS() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return *client
}

func peticionGET() {
	client := ignorarHTTPS()
	resp, err := client.Get(urlServer)
	chkError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	chkError(err)
	fmt.Println(string(body))
}

// función para codificar de []bytes a string (Base64)
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// función para decodificar de string a []bytes (Base64)
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chkError(err)                                // comprobamos el error
	return b                                     // devolvemos los datos originales
}

func login() bool {
	fmt.Printf("\n__Login__\n")

	//Pedir datos
	var user string
	fmt.Printf("Usuario: ")
	fmt.Scanf("%s\n", &user)

	var password string
	fmt.Printf("Contraseña: ")
	fmt.Scanf("%s\n", &password)

	hasher := md5.New()
	hasher.Write([]byte(password))
	password = hex.EncodeToString(hasher.Sum(nil))

	//serializar a JSON
	m := userRes{user, password}
	loginJSON, err := json.Marshal(m)
	chkError(err)
	correct := loginPost(loginJSON)

	if correct {
		fmt.Printf("Bienvenido!\n\n")
		menuLogueado(user)
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func loginPost(js []byte) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "login")
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data) // enviamos por POST
	chkError(err)

	var respJS resp
	//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	json.NewDecoder(r.Body).Decode(&respJS)
	if respJS.Ok {
		return true
	}
	return false
}

func registroPost(js []byte) bool {
	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "registro")
	data.Set("mensaje", encode64(js))
	client.PostForm(urlServer, data)
	fmt.Println("a")
	return true
}

func consultarCuentas(boss string) bool { //boss es el nombre del usuario logueado
	fmt.Printf("\n__Tus cuentas__\n")

	//serializar a JSON
	m := cuentaRes{boss, "", "", ""}
	cuentaJSON, err := json.Marshal(m)
	chkError(err)
	correct := consultarCuentasPost(cuentaJSON)

	return correct
}

func consultarCuentasPost(js []byte) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "consultarCuentas")
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data) // enviamos por POST
	chkError(err)

	var respJS respJSON
	//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	json.NewDecoder(r.Body).Decode(&respJS)
	fmt.Println(imprimirConsulta(respJS.Info))

	if imprimirConsulta(respJS.Info) == "No hay ninguna cuenta.\n" {
		return false
	}
	return true
}

func imprimirConsulta(info map[string]datos) string {
	var s string
	if len(info) == 0 {
		s = "No hay ninguna cuenta.\n"
	} else {

		for key, val := range info {
			s += fmt.Sprintf("#%s:\n\tUsuario: %s\n\tContraseña: %s\n", key, val.User, val.Pass)
		}
	}
	return s
}

func añadirCuenta(boss string) bool { //boss es el nombre del usuario logueado
	fmt.Printf("\n__Añadir nueva cuenta__\n")

	//Pedir datos
	var servicio string
	fmt.Printf("Nuevo servicio: ")
	fmt.Scanf("%s\n", &servicio)

	//Pedir datos
	var user string
	fmt.Printf("Nuevo nombre de usuario: ")
	fmt.Scanf("%s\n", &user)

	var password string
	fmt.Printf("Nueva contraseña: ")
	fmt.Scanf("%s\n", &password)

	//serializar a JSON
	m := cuentaRes{boss, servicio, user, password}
	cuentaJSON, err := json.Marshal(m)
	chkError(err)
	correct := añadirCuentaPost(cuentaJSON)

	if correct {
		fmt.Printf("Añadida correctamente!\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func añadirCuentaPost(js []byte) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "añadirCuenta")
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data) // enviamos por POST
	chkError(err)

	var respJS resp
	//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	json.NewDecoder(r.Body).Decode(&respJS)
	if respJS.Ok {
		return true
	}
	return false
}

func eliminarCuenta(boss string) bool { //boss es el nombre del usuario logueado

	if consultarCuentas(boss) == false {
		return false
	}

	fmt.Printf("\n__Eliminar cuenta__\n")

	//Pedir datos
	var servicio string
	fmt.Printf("Selecciona el servicio: ")
	fmt.Scanf("%s\n", &servicio)

	//serializar a JSON
	m := cuentaRes{boss, servicio, "", ""}
	cuentaJSON, err := json.Marshal(m)
	chkError(err)
	correct := eliminarCuentaPost(cuentaJSON)

	if correct {
		fmt.Printf("\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func eliminarCuentaPost(js []byte) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "eliminarCuenta")
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data) // enviamos por POST
	chkError(err)

	var respJS resp
	//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	json.NewDecoder(r.Body).Decode(&respJS)
	if respJS.Ok {
		return true
	}
	return false
}

func main() {

	var opcion = menu()
	switch opcion {
	case 1:
		login()
	case 2:
		registro()
	case 3:
		break
	default:
		break
	}
}

func registro() bool {
	var user, passwd, tempPasswd, mail string
	var correct bool
	//Pedimos el nombre de usuario
	fmt.Println("Introduce tu nombre de usuario")
	n, err := fmt.Scanf("%s\n", &user)
	if err != nil {
		fmt.Println(n, err)
	}
	//Pedimos la contraseña
	for {
		fmt.Println("Introduce tu contraseña")
		n, err = fmt.Scanf("%s\n", &passwd)
		if err != nil {
			fmt.Println(n, err)
		}

		//Volvemos a pedir la contraseña
		fmt.Println("Vuelve a introducir tu contraseña")
		n, err = fmt.Scanf("%s\n", &tempPasswd)
		if err != nil {
			fmt.Println(n, err)
		}
		if passwd == tempPasswd {
			break
		} else {
			fmt.Println("Las contraseñas no coinciden")
		}
	}
	//Pedimos el email
	fmt.Printf("Introduce tu email\n")
	n, err = fmt.Scanf("%s\n", &mail)
	if err != nil {
		fmt.Println(n, err)
	}
	//Generamos el hash a partir de la contraseña
	hasher := md5.New()
	hasher.Write([]byte(passwd))
	passwd = hex.EncodeToString(hasher.Sum(nil))

	//Ahora almacenamos el usuario en formato Json
	newUser := structUser{user, passwd, mail}
	b, error := json.Marshal(&newUser)
	if err != nil {
		fmt.Println(error)
	}
	correct = registroPost(b)
	if correct {
		fmt.Printf("Registrado correctamente\n")
	}
	return correct
}

func generarPassword() {
	var str_size int
	fmt.Printf("Inserta longitud del Password: ")
	fmt.Scanf("%d\n", &str_size)
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	fmt.Println(string(bytes))
}

func menuLogueado(username string) {
	var opcion int
	fmt.Printf("----------Bienvenido %s-------", username)
	fmt.Println("-------------------------")
	fmt.Printf("1 - Consultar cuentas\n")
	fmt.Printf("2 - Agregar cuenta\n")
	fmt.Printf("3 - Eliminar cuenta\n")
	fmt.Printf("4 - Generar Password\n")
	fmt.Printf("5 - Salir\n")
	fmt.Printf("Opción: ")
	fmt.Scanf("%d\n", &opcion)
	switch opcion {
	case 1:
		consultarCuentas(username)
	case 2:
		añadirCuenta(username)
	case 3:
		eliminarCuenta(username)
	case 4:
		generarPassword()
	default:
		break
	}
}
