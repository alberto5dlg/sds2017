package main

import (
	"fmt"
	"strings"
)

var allrunes = map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
	"Ñ": 14, "O": 15, "P": 16, "Q": 17, "R": 18, "S": 19, "T": 20, "U": 21, "V": 22, "W": 23, "X": 24, "Y": 25, "Z": 26}

var invallrunes = map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M", 13: "N",
	14: "Ñ", 15: "O", 16: "P", 17: "Q", 18: "R", 19: "S", 20: "T", 21: "U", 22: "V", 23: "W", 24: "X", 25: "Y", 26: "Z"}

func main() {

	var codigo string
	var codigoCifrar []string
	var codigoCifrado []string
	fmt.Printf("Introduce el codigo para cifrar\n")
	n, err := fmt.Scanf("%s\n", &codigo)
	if err != nil {
		fmt.Println(n, err)
	}

	var desplazamiento int
	fmt.Printf("Introduce el desplazamiento para cifrar\n")
	n, err = fmt.Scanf("%d\n", &desplazamiento)
	if err != nil {
		fmt.Println(n, err)
	}
	codigo = strings.ToUpper(codigo)
	codigoCifrar = strings.SplitAfter(codigo, "")
	codigoCifrado = strings.SplitAfter(codigo, "")
	for i := range codigoCifrar {
		var posicion int = allrunes[codigoCifrar[i]]
		posicion = posicion + desplazamiento
		for posicion > 26 {
			posicion = posicion % 26
		}
		codigoCifrado[i] = invallrunes[posicion]
	}
	var codigoCifradofinal string = strings.Join(codigoCifrado, "")
	fmt.Printf("El codigo 22 es %s y el desplazamiento %d y el codigo cifrado es %s", codigo, desplazamiento, codigoCifradofinal)

}
