package main

import (
	"fmt"
	"os"
	"strings"
)

// Delimitador que separa no terminales y terminales de las producciones
const PRODUCTION_DELIMITER = " -> "

type FIRSTS map[string][]string
type FOLLOWS map[string][]string

var FIRSTS_CACHE_STATE = make(FIRSTS)
var FOLLOWS_CACHE_STATE = make(FOLLOWS)

// Estructura de analisis
type AnalyzerOutput struct {
	terminals, non_terminals []string
}

// Funcion que dado una lista de cadenas izquierda y derecha
// regrese una lista de cadenas representando las terminales de
// un conjunto de producciones.
func GetTerminals(left []string, right []string) []string {
	var terminals []string

	right_size := len(right)

	// Para cada produccion, conseguir los tokens separando las
	// cadenas con el delimitador espacio.
	for i := 0; i < right_size; i++ {
		tokens := strings.Split(right[i], " ")
		tokens_size := len(tokens)
		for j := 0; j < tokens_size; j++ {
			token := tokens[j]
			// Agregar a la lista los tokens que no sean no terminales, y que no sean epsilon
			if !Contains(left, token) && token != "'" && token != "" {
				terminals = append(terminals, token)
			}
		}
	}

	// Regresar terminales eliminando duplicados
	return RemoveDuplicates(terminals)
}

// Funcion que regresa no terminales dado una lista de cadenas
// eliminando los duplicados
func GetNonTerminals(left []string) []string {
	return RemoveDuplicates(left)
}

// Funcion que analiza las producciones para regresar una instancia de AnalyzerOutput
// dada una lista de cadenas
func Analyze(lines []string) AnalyzerOutput {
	var left []string
	var right []string

	size := len(lines)

	// Para cada produccion, obtener izquierda y derecha con el delimitador.
	for i := 0; i < size; i++ {
		splitted := strings.Split(lines[i], PRODUCTION_DELIMITER)
		left = append(left, splitted[0])
		right = append(right, splitted[1])
	}

	// Llamadas a funciones para regresar terminales y no terminales
	nonTerminals := GetNonTerminals(left)
	terminals := GetTerminals(left, right)

	return AnalyzerOutput{terminals, nonTerminals}
}

func FindFirst(lines []string, non_terminal string, analyzer AnalyzerOutput) []string {

	cache, found := FIRSTS_CACHE_STATE[non_terminal]
	if found {
		return cache
	}

	var result []string

	size := len(lines)

	for i := 0; i < size; i++ {
		splitted := strings.Split(lines[i], PRODUCTION_DELIMITER)
		left := splitted[0]

		if left != non_terminal {
			continue
		}

		right := splitted[1]

		first_right := strings.Split(right, " ")[0]

		if non_terminal == first_right {
			//TODO: IMprove error handling

			// Recursion found. Not LL1.
			fmt.Fprintf(os.Stderr, "Not valid LL(1) format.\n")
			os.Exit(1)
		}

		// If epsilon
		if first_right == "'" {
			continue
		}

		// If is a non terminal, find firsts of that non terminal
		if Contains(analyzer.non_terminals, first_right) {
			found := FindFirst(lines, first_right, analyzer)
			result = append(result, found...)
		} else {
			result = append(result, first_right)
		}

	}
	// Save to cache
	FIRSTS_CACHE_STATE[non_terminal] = result
	return result
}

func GetFirsts(lines []string, analyzer AnalyzerOutput) FIRSTS {
	size := len(analyzer.non_terminals)

	for i := 0; i < size; i++ {
		non_terminal := analyzer.non_terminals[i]
		FindFirst(lines, non_terminal, analyzer)
	}

	return FIRSTS_CACHE_STATE
}
