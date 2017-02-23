package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	var longMatriz int
	var frase string

	fmt.Printf("Introduce la frase:\n ")
	in := bufio.NewReader(os.Stdin)
	lineBuf, err := in.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error en ReadLine: ", err)
		os.Exit(2)
	}
	frase = string(lineBuf)

	fmt.Println("Introduce la longitud de la escítala")
	n, err := fmt.Scanf("%d\n", &longMatriz)
	if err != nil {
		fmt.Println(n, err)
	}
	fmt.Printf("La longitud de la fila matriz es %d\n", longMatriz)
	fmt.Printf("La frase para cifrar es %s", frase)
	var codigoCifrar []string
	codigoCifrar = strings.SplitAfter(frase, "")
	var altMatriz int = len(frase) / longMatriz
	var matriz [longMatriz][altMatriz]string

	frase = strings.ToUpper(frase)
	var contFila int = 0
	var contColumn int = 0
	for i := range codigoCifrar {

		if i < longMatriz {
			matriz[contFila][contColumn] = codigoCifrar[i]
			contFila = contFila + 1
		}
		if i >= longMatriz {
			contFila = 0
			contColumn++
		}
	}

}
