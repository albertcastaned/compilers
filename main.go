// Alberto Castañeda Arana
// A01250647

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var productions_amount int

	// Se inicia instancia de scanner para contemplar espacios vacios
	scanner := bufio.NewScanner(os.Stdin)

	// Escanear numero de producciones a ingresar
	fmt.Scanf("%d", &productions_amount)
	var lines []string

	// Para cada numero especificado, leer las entradas de producciones
	for i := 0; i < productions_amount; i++ {
		scanner.Scan()
		lines = append(lines, scanner.Text())
	}

	// Llamar a analizar las producciones
	output := Analyze(lines)

	// Imprimir terminales y no terminales
	fmt.Printf(
		"Terminal: %s\nNon terminal: %s\n",
		strings.Join(output.terminals, ", "),
		strings.Join(output.non_terminals, ", "),
	)

	size := len(output.non_terminals)

	firsts_state := GetFirsts(lines, output)
	for i := 0; i < size; i++ {
		non_terminal := output.non_terminals[i]
		fmt.Printf("%s => FIRST = {%s}\n", non_terminal, strings.Join(firsts_state[non_terminal], ", "))
	}
}
