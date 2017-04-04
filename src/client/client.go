package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var urlServer = "https://127.0.0.1:8081"

type userRes struct {
	Method   string
	User     string
	Password string
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
	var method = "login"
	m := userRes{method, user, password}
	mJSON, err := json.Marshal(m)
	chkError(err)
	os.Stdout.Write(mJSON)
	fmt.Printf("\n")

	//Encriptar la información
	//PASARLO A BASE64 ANTES DE ENVIARLO PARA QUE NO DE PROBLEMAS EL TIPO BYTE

	if correct == false {
		fmt.Printf("User or password Error\n\n")
	} else {
		fmt.Printf("Welcome!\n\n")
	}
	return correct
}

func main() {

	var opcion = menu()
	switch opcion {
	case 1:
		login()
	case 3:
		break
	default:
		break
	}
}
