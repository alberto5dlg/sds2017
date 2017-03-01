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

func main() {
	var url = urlServer + ":" + portServer
	// connect to this socket
	conn, _ := net.Dial(typeConexion, url)
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
