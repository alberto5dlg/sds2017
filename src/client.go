package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var urlServer = "https://127.0.0.1:8081"
var keyCifrado []byte
var token string

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

	//Generamos el hash a partir de la contraseña
	hasher := sha512.Sum512([]byte(password))
	keyCifrado = hasher[32:64] // Utilizaremos la segunda mitad como key para el cifrado
	password = encode64(hasher[:])

	//serializar a JSON
	m := logueado{user, password}
	loginJSON, err := json.Marshal(m)
	chkError(err)
	correct := metodoPostLoginRegistro(loginJSON, "login")

	if correct {
		if compDA(user) {
			fmt.Printf("Bienvenido!\n\n")
			menuLogueado(user)
		} else {
			fmt.Printf("Error en Doble Autenticacion!\n\n")
		}
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func compDA(user string) bool {
	var clave string
	fmt.Printf("Introduce la clave enviada al Mail: ")
	fmt.Scanf("%s\n", &clave)
	da := struAuth{clave}
	daJSON, err := json.Marshal(da)
	chkError(err)
	correct := metodoPost(daJSON, "dobleAuth", user)
	return correct
}

func consultarCuentas(boss string) bool { //boss es el nombre del usuario logueado
	fmt.Printf("\n__Tus cuentas__\n")
	//serializar a JSON
	m := cuenta{boss, "", "", ""}
	cuentaJSON, err := json.Marshal(m)
	chkError(err)
	correct := consultarCuentasPost(cuentaJSON, boss)

	return correct
}

func consultarCuentasPost(js []byte, username string) bool {
	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "consultarCuentas")
	data.Set("username", username)
	data.Set("token", token)
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
func descifrarPassword(tempPass string) []byte {
	password := decode64(tempPass)
	ciphertext, err := decrypt(password, keyCifrado)
	chkError(err)
	return ciphertext
}
func imprimirConsulta(info map[string]datos) string {
	var s string
	if len(info) == 0 {
		s = "No hay ninguna cuenta.\n"
	} else {
		for key, val := range info {
			var password = descifrarPassword(val.Pass)
			s += fmt.Sprintf("#%s:\n\tUsuario: %s\n\tContraseña: %s\n", key, val.User, password)
		}
	}
	return s
}
func decrypt(password []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(password) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, password := password[:nonceSize], password[nonceSize:]
	return gcm.Open(nil, nonce, password, nil)
}
func encrypt(password []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, password, nil), nil
}

func cifrarPassword(tempPass string) string {
	password := []byte(tempPass)
	ciphertext, err := encrypt(password, keyCifrado)
	chkError(err)
	var result = encode64(ciphertext)
	return result
}

func anyadirCuenta(boss string) bool { //boss es el nombre del usuario logueado
	fmt.Printf("\n__Añadir nueva cuenta__\n")
	var nCuenta cuenta
	nCuenta.Boss = boss
	fmt.Printf("Nuevo servicio: ")
	fmt.Scanf("%s\n", &nCuenta.Servicio)
	fmt.Printf("Nuevo nombre de usuario: ")
	fmt.Scanf("%s\n", &nCuenta.User)
	fmt.Printf("Nueva contraseña: ")
	fmt.Scanf("%s\n", &nCuenta.Password)
	nCuenta.Password = cifrarPassword(nCuenta.Password)

	//serializar a JSON
	cuentaJSON, err := json.Marshal(nCuenta)
	chkError(err)
	correct := metodoPost(cuentaJSON, "añadirCuenta", boss)
	if correct {
		fmt.Printf("Añadida correctamente!\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
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
	m := cuenta{boss, servicio, "", ""}
	cuentaJSON, err := json.Marshal(m)
	chkError(err)
	correct := eliminarCuentaPost(cuentaJSON, boss)

	if correct {
		fmt.Printf("\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func eliminarCuentaPost(js []byte, username string) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "eliminarCuenta")
	data.Set("username", username)
	data.Set("token", token)
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
	fmt.Printf("Introduce tu nombre de usuario: ")
	n, err := fmt.Scanf("%s\n", &user)
	if err != nil {
		fmt.Println(n, err)
	}
	//Pedimos la contraseña
	for {
		fmt.Printf("Introduce tu contraseña: ")
		n, err = fmt.Scanf("%s\n", &passwd)
		if err != nil {
			fmt.Println(n, err)
		}

		//Volvemos a pedir la contraseña
		fmt.Printf("Vuelve a introducir tu contraseña: ")
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
	fmt.Printf("Introduce tu email: ")
	n, err = fmt.Scanf("%s\n", &mail)
	chkError(err)
	//Generamos el hash a partir de la contraseña
	hasher := sha512.Sum512([]byte(passwd))
	keyCifrado = hasher[32:64] // Utilizaremos la segunda mitad como key para el cifrado
	passwd = encode64(hasher[:])

	//Ahora almacenamos el usuario en formato Json
	newUser := registrarse{user, passwd, mail}
	b, err := json.Marshal(&newUser)
	chkError(err)
	correct = metodoPostLoginRegistro(b, "registro")
	if correct {
		fmt.Printf("Registrado correctamente\n")
	}
	menuLogueado(newUser.User)
	return correct
}

func generarPassword() {
	var strSize int
	fmt.Printf("Inserta longitud del Password: ")
	fmt.Scanf("%d\n", &strSize)
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	fmt.Println(string(bytes))
}

func anyadirTarjeta(username string) bool {
	fmt.Printf("\n__Añadir nueva Tarjeta__\n")
	var nCard nTarjeta
	nCard.Username = username
	fmt.Printf("Entidad: ")
	fmt.Scanf("%s\n", &nCard.Entidad)
	fmt.Printf("Numero de Tarjeta: ")
	fmt.Scanf("%s\n", &nCard.NTarjeta)
	fmt.Printf("Codigo de Seguridad: ")
	fmt.Scanf("%s\n", &nCard.CodSeg)
	fmt.Printf("Fecha de tarjeta: ")
	fmt.Scanf("%s\n", &nCard.Fecha)
	nCard.CodSeg = cifrarPassword(nCard.CodSeg)
	nCard.Fecha = cifrarPassword(nCard.Fecha)
	nCard.NTarjeta = cifrarPassword(nCard.NTarjeta)

	cuentaJSON, err := json.Marshal(nCard)
	chkError(err)
	correct := metodoPost(cuentaJSON, "añadirTarjeta", username)
	if correct {
		fmt.Printf("Añadida correctamente!\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func anyadirNotas(username string) bool {
	fmt.Printf("__ Añadir nueva nota __\n")
	reader := bufio.NewReader(os.Stdin)
	var nNote nNotas
	nNote.Username = username
	fmt.Printf("Titulo: \n")
	text, _ := reader.ReadString('\n')
	nNote.Titulo = text
	fmt.Printf("Texto: \n")
	text, _ = reader.ReadString('\n')
	nNote.Cuerpo = text

	//nCuenta.Password = cifrarPassword(nCuenta.Password)
	nNote.Titulo = cifrarPassword(nNote.Titulo)
	nNote.Cuerpo = cifrarPassword(nNote.Cuerpo)

	cuentaJSON, err := json.Marshal(nNote)
	chkError(err)
	correct := metodoPost(cuentaJSON, "añadirNota", username)
	if correct {
		fmt.Printf("Nueva Nota!\n\n")
	} else {
		fmt.Printf("Error!\n\n")
	}
	return correct
}

func metodoPost(js []byte, comando string, username string) bool {
	client := ignorarHTTPS()
	data := url.Values{}
	data.Set("cmd", comando)
	data.Set("username", username)
	data.Set("token", token)
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data)
	chkError(err)
	var respJS resp
	json.NewDecoder(r.Body).Decode(&respJS)

	if respJS.Ok {
		if comando == "login" || comando == "registro" {
			token = respJS.Msg
		}
		return true
	}
	return false
}

func metodoPostLoginRegistro(js []byte, comando string) bool {
	client := ignorarHTTPS()
	data := url.Values{}
	data.Set("cmd", comando)
	data.Set("mensaje", encode64(js))
	r, err := client.PostForm(urlServer, data)
	chkError(err)
	var respJS resp
	json.NewDecoder(r.Body).Decode(&respJS)

	if respJS.Ok {
		if comando == "login" || comando == "registro" {
			token = respJS.Msg
		}
		return true
	}
	return false
}

func menuLogueado(username string) {
	var opcion int
	for opcion != 7 {
		fmt.Printf("----------Bienvenido %s-------", username)
		fmt.Println("-------------------------")
		fmt.Printf("1 - Consultar cuentas\n")
		fmt.Printf("2 - Agregar cuenta\n")
		fmt.Printf("3 - Agregar Tarjeta\n")
		fmt.Printf("4 - Agregar Nota\n")
		fmt.Printf("5 - Eliminar cuenta\n")
		fmt.Printf("6 - Generar Password\n")
		fmt.Printf("7 - Salir\n")
		fmt.Printf("Opción: ")
		fmt.Scanf("%d\n", &opcion)
		switch opcion {
		case 1:
			consultarCuentas(username)
		case 2:
			anyadirCuenta(username)
		case 3:
			anyadirTarjeta(username)
		case 4:
			anyadirNotas(username)
		case 5:
			eliminarCuenta(username)
		case 6:
			generarPassword()
		case 7:
			break
		default:
			fmt.Println("Opcion Incorrecta !!")
		}
	}
}