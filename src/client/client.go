package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

var urlServer = "https://127.0.0.1:8081"

func menu() int {
	var opcion int
	for opcion <= 0 || opcion >= 4 {
		fmt.Printf("Aplicación SDS Seguridad\n")
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("1 - Login\n")
		fmt.Printf("2 - Registro\n")
		fmt.Printf("3 - Salir\n")
		fmt.Printf("Opción: ")
		n, err := fmt.Scanf("%d\n", &opcion)
		if err != nil {
			fmt.Println(n, err)
		}
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
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func main() {

	opc := menu()
	switch opc {
	case 1:
		peticionGET()
		break
	case 3:
		println("Adios")
	}
}
