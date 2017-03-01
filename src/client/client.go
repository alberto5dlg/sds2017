package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//
var urlServer = "127.0.0.1"
var portServer = "8081"
var typeConexion = "tcp"

func menu() int {
	var opcion int = 0
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
func main() {
	var url = urlServer + ":" + portServer
	// connect to this socket
	conn, _ := net.Dial(typeConexion, url)
	var opcion int = menu()
	fmt.Printf("%d\n", opcion)
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
