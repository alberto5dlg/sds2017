package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var urlServer = "https://127.0.0.1:8081"

type userRes struct {
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

	var user string
	fmt.Printf("User: ")
	err := fmt.Scanf("%s\n", &user)
	chkError(err)

	var password string
	fmt.Printf("Password: ")
	err = fmt.Scanf("%s\n", &password)
	chkError(err)

	//llamada al servidor / comprobar paswword
	m := userRes{user, password}
	mJSON, err := json.Marshal(m)
	chkError(err)
	os.Stdout.Write(mJSON)
	fmt.Printf("\n")

	if correct == false {
		fmt.Printf("User or password Error\n\n")
	} else {
		fmt.Printf("Welcome!\n\n")
	}
	return correct
}

func main() {
	var url = urlServer + ":" + portServer
	// connect to this socket
	conn, _ := net.Dial(typeConexion, url)
	var opcion = menu()
	if opcion == 1 {
		login()
	}
	for {

		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: " + message)
	}
}
