// Alberto Castañeda Arana
// A01250647

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var productions_amount, inputs_amount int
	// Se inicia instancia de scanner para contemplar espacios vacios
	scanner := bufio.NewScanner(os.Stdin)

	// Escanear numero de producciones a ingresar
	fmt.Scanf("%d %d\n", &productions_amount, &inputs_amount)
	var lines, inputs []string

	// Para cada numero especificado, leer las entradas de producciones
	for i := 0; i < productions_amount; i++ {
		scanner.Scan()
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < inputs_amount; i++ {
		scanner.Scan()
		inputs = append(inputs, scanner.Text())
	}

	// Llamar a analizar las producciones
	output := Analyze(lines)

	// Imprimir terminales y no terminales
	// fmt.Printf(
	// 	"Terminal: %s\nNon terminal: %s\n",
	// 	strings.Join(output.terminals, ", "),
	// 	strings.Join(output.non_terminals, ", "),
	// )

	// size := len(output.non_terminals)

	// for i := 0; i < size; i++ {
	// 	non_terminal := output.non_terminals[i]
	// 	firsts := FindFirst(lines, non_terminal, output)
	// 	follows := FindFollow(lines, non_terminal, output, "")

	// 	fmt.Printf("%s => FIRST = {%s}, FOLLOW = {%s}\n", non_terminal, strings.Join(firsts, ", "), strings.Join(follows, ", "))
	// }

	var ll1_valid bool
	ll1_valid = IsLL1Valid(lines, output)

	if ll1_valid == false {
		fmt.Printf("Grammar is not LL(1)!\n")
		os.Exit(-1)
	}

	GetLL1Table(lines, output)

	for _, input := range inputs {
		var result string
		if CheckValidInput(lines, input) {
			result = "YES"
		} else {
			result = "NO"
		}
		fmt.Printf("%s - ACCEPTED? %s\n", input, result)
	}
}
