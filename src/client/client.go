package main

import (
	"crypto/md5"
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

type userRes struct {
	User     string
	Password string
}

type structUser struct {
	User     string
	Password string
	Email    string
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
	var correct = false
	fmt.Printf("\n__Login__\n")

	//Pedir datos
	var user string
	fmt.Printf("User: ")
	fmt.Scanf("%s\n", &user)

	var password string
	fmt.Printf("Password: ")
	fmt.Scanf("%s\n", &password)

	//serializar a JSON
	m := userRes{user, password}
	loginJSON, err := json.Marshal(m)
	chkError(err)
	loginPost(loginJSON)

	//Encriptar la información
	//PASARLO A BASE64 ANTES DE ENVIARLO PARA QUE NO DE PROBLEMAS EL TIPO BYTE

	if correct == false {
		fmt.Printf("User or password Error\n\n")
	} else {
		fmt.Printf("Welcome!\n\n")
	}
	return correct
}

func loginPost(js []byte) bool {

	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "login")
	data.Set("mensaje", encode64(js))
	client.PostForm(urlServer, data) // enviamos por POST
	fmt.Printf("enviado!\n\n")
	//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	//fmt.Println()

	return false
}
func registroPost(js []byte) bool {
	client := ignorarHTTPS()

	data := url.Values{}
	data.Set("cmd", "registro")
	data.Set("mensaje", encode64(js))
	client.PostForm(urlServer, data)
	return true
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
	var correct bool = false
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
	hex.EncodeToString(hasher.Sum(nil))

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
