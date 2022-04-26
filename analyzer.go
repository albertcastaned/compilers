package main

import (
	"strings"
)

// Delimitador que separa no terminales y terminales de las producciones
const PRODUCTION_DELIMITER = " -> "

type FIRSTS map[string][]string
type FOLLOWS map[string][]string

var FIRSTS_CACHE_STATE = make(FIRSTS)
var FOLLOWS_CACHE_STATE = make(FOLLOWS)

var LL1_VALID = true

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
		a, b := SplitProduction(lines[i])
		left = append(left, a)
		right = append(right, b)
	}

	// Llamadas a funciones para regresar terminales y no terminales
	non_terminals := GetNonTerminals(left)
	terminals := GetTerminals(left, right)

	return AnalyzerOutput{terminals, non_terminals}
}

func FindFirst(productions []string, value string, analyzer AnalyzerOutput) []string {

	cache, found := FIRSTS_CACHE_STATE[value]
	if found {
		return cache
	}

	var result []string

	// Value is a terminal, first of a terminal is the terminal itself..
	if Contains(analyzer.terminals, value) {
		var result []string = []string{value}
		return result
	}

	size := len(productions)

	for i := 0; i < size; i++ {
		left, right := SplitProduction(productions[i])

		if left != value {
			continue
		}

		first_right := strings.Split(right, " ")[0]

		if value == first_right {
			LL1_VALID = false
			continue
		}

		// If epsilon
		if first_right == "'" {
			result = append(result, "' '")
			continue
		}

		// If is a non terminal, find firsts of that non terminal
		if Contains(analyzer.non_terminals, first_right) {
			found := FindFirst(productions, first_right, analyzer)
			result = append(result, found...)
		} else {
			result = append(result, first_right)
		}

	}
	// Save to cache
	FIRSTS_CACHE_STATE[value] = RemoveDuplicates(result)
	return result
}

func GetFirsts(productions []string, analyzer AnalyzerOutput) (FIRSTS, bool) {
	size := len(analyzer.non_terminals)

	for i := 0; i < size; i++ {
		non_terminal := analyzer.non_terminals[i]
		FindFirst(productions, non_terminal, analyzer)
	}

	return FIRSTS_CACHE_STATE, LL1_VALID
}

func FindFollow(productions []string, value string, analyzer AnalyzerOutput) []string {
	cache, found := FOLLOWS_CACHE_STATE[value]
	if found {
		return cache
	}

	var follows []string

	is_start_symbol := strings.Split(productions[0], PRODUCTION_DELIMITER)[0] == value

	// First rule
	if is_start_symbol {
		follows = append(follows, "$")
	}

	for _, production := range productions {
		left, right := SplitProduction(production)

		tokens := strings.Split(right, " ")
		size_tokens := len(tokens)

		//TODO: CHECK RECURSION
		for index, token := range tokens {
			third_rule := false

			if token == value {

				// Is last element
				if index+1 == size_tokens {
					// Rule 3: B -> aA where A is our desired nonterminal
					third_rule = true
				} else {
					// Rule 3: B -> aAB only when B -> epsilon.
					if tokens[index+1] == "'" {
						third_rule = true
					} else {
						// Rule 2: B -> aAB -> FOLLOW(A) = FIRST(B)
						// Get firsts of next token, previous if statement protects from out of array bounds
						follows = append(follows, FindFirst(productions, tokens[index+1], analyzer)...)
					}
				}

				// Rule 3: FOLLOW(A) = FOLLOW(B) where A would be our left and B our current token
				if third_rule {

					// Protect against recursion
					if left == value {
						LL1_VALID = false
						continue
					}
					follows = append(follows, FindFollow(productions, left, analyzer)...)
				}
			}
		}

	}

	FOLLOWS_CACHE_STATE[value] = RemoveDuplicates(follows)
	return follows
}

func GetFollows(productions []string, analyzer AnalyzerOutput) (FOLLOWS, bool) {
	size := len(analyzer.non_terminals)
	for i := 0; i < size; i++ {
		non_terminal := analyzer.non_terminals[i]
		FindFollow(productions, non_terminal, analyzer)
	}

	return FOLLOWS_CACHE_STATE, LL1_VALID
}
