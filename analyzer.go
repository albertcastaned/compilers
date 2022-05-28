package main

import (
	"fmt"
	"os"
	"strings"
)

// Delimitador que separa no terminales y terminales de las producciones
const PRODUCTION_DELIMITER = " -> "

// Inicializamos cache de FIRSTS y de FOLLOWS
type FIRSTS map[string][]string
type FOLLOWS map[string][]string

var FIRSTS_CACHE_STATE = make(FIRSTS)
var FOLLOWS_CACHE_STATE = make(FOLLOWS)

type LL_TABLE_VALUE map[string]string

type LL1_TABLE map[string]LL_TABLE_VALUE

var LL1_TABLE_STATE = make(LL1_TABLE)

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

// Funcion para regresar los FIRSTS de una produccion especifica
func FindProductionFirst(productions []string, index int, analyzer AnalyzerOutput) []string {
	var result []string

	// Obtenemos izquierdo y derecho de produccion
	left, right := SplitProduction(productions[index])

	first_right := strings.Split(right, " ")[0]

	// Recuperamos de cache los firsts si ya han sido calculados antes
	cache, found := FIRSTS_CACHE_STATE[first_right]
	if found {
		return cache
	}

	// Manejamos casos de recursion infinito
	if left == first_right {
		return result
	}

	// Si encontramos un epsilon, lo agregamos en formato correcto
	if first_right == "'" {
		result = append(result, "' '")
		return result
	}

	// Si es un non terminal, llamar recursivamente a buscar los firsts de ese non terminal
	if Contains(analyzer.non_terminals, first_right) {
		found := FindFirst(productions, first_right, analyzer)
		result = append(result, found...)
	} else {
		// Si es un terminal, agregamos a la lista de firsts
		result = append(result, first_right)
	}

	return result
}

// Funcion para encontrar todos los FIRSTS dado un no terminal
func FindFirst(productions []string, value string, analyzer AnalyzerOutput) []string {

	// Recuperamos de cache los firsts si ya han sido calculados antes
	cache, found := FIRSTS_CACHE_STATE[value]
	if found {
		return cache
	}

	var result []string

	// Verificamos que el valor es un terminal y si es asi regresar al terminal en si
	if Contains(analyzer.terminals, value) {
		var result []string = []string{value}
		return result
	}

	size := len(productions)

	// Para cada produccion, separar izquierda y derecha
	for i := 0; i < size; i++ {
		left, _ := SplitProduction(productions[i])

		// Si la produccion empieza con el no terminal deseado, buscamos los firsts de esta produccion.
		if left == value {
			result = append(result, FindProductionFirst(productions, i, analyzer)...)
		}
	}

	// Eliminamos posibles duplicados de resultado
	result = RemoveDuplicates(result)

	// Guardamos a cache o dp para posterior recuperacion.
	FIRSTS_CACHE_STATE[value] = result
	return result
}

// Funcion para encontrar los follows dado una lista de producciones y un valor
func FindFollow(productions []string, value string, analyzer AnalyzerOutput, lastCall string) []string {

	// Recuperamos de cache los follows si ya han sido calculados antes
	cache, found := FOLLOWS_CACHE_STATE[value]
	if found && len(cache) > 0 {
		return cache
	}

	var follows []string

	// Si es el no terminal inicial, o el izquierdo de la primera produccion, se activa la primera regla
	is_start_symbol := strings.Split(productions[0], PRODUCTION_DELIMITER)[0] == value

	// Rule 1: Add $ to start symbol
	if is_start_symbol {
		follows = append(follows, "$")
	}

	// Para cada produccion:
	for _, production := range productions {
		left, right := SplitProduction(production)

		// Obteenmos tokens del lado derecho de la produccion
		tokens := strings.Split(right, " ")
		size_tokens := len(tokens)

		// Para cada token de la produccion
		for index, token := range tokens {
			third_rule := false

			// Si el token is igual al no terminal buscado:
			if token == value {

				// Si es el ultimo elemento
				if index+1 == size_tokens {
					// Rule 3: B -> aA where A is our desired nonterminal
					third_rule = true
				} else {
					// Rule 3: B -> aAB only when B -> epsilon.

					firsts_of_right := Contains(FindFirst(productions, tokens[index+1], analyzer), "' '")
					if firsts_of_right {
						third_rule = true
					}
					// Rule 2: B -> aAB -> FOLLOW(A) = FIRST(B)
					// Obtenemos los firsts del siguiente token.
					follows = append(follows, FindFirst(productions, tokens[index+1], analyzer)...)
				}

				// Rule 3: FOLLOW(A) = FOLLOW(B) where A would be our left and B our current token
				if third_rule {
					// Proteccion contra posible recursion infinita
					if lastCall == left {
						continue
					}

					// Con la tercer regla, agregamos y buscamos los follows de la izquierda
					follows = append(follows, FindFollow(productions, left, analyzer, value)...)
				}
			}
		}

	}

	// Eliminamos los epsilons encontrados y los duplicados
	follows = RemoveEpsilons(RemoveDuplicates(follows))

	// Guardamos a cache de follows
	FOLLOWS_CACHE_STATE[value] = follows

	return follows
}

// Funcion para verificar que dado una lista de producciones, la gramatica es LL(1) valida.
func IsLL1Valid(productions []string, analyzer AnalyzerOutput) bool {
	non_terminals := analyzer.non_terminals

	// Para cada no terminal
	for _, non_terminal := range non_terminals {
		var found []string
		var indexes []int

		// Para cada produccion, contamos cuantas hay que empiezan con el no terminal en iteracion.
		for index, production := range productions {
			left, _ := SplitProduction(production)
			if left == non_terminal {
				found = append(found, production)
				indexes = append(indexes, index)
			}
		}

		// Si solo se tiene una produccion para este no terminal, no hay ninguna regla que checar.
		if len(found) <= 1 {
			continue
		}

		// Creamos posibles combinaciones de pares de este no terminal en las producciones
		combinations := CreateCombinations(len(indexes))

		for _, combination := range combinations {
			// Obtenemos firsts del par
			firsts_a := FindProductionFirst(productions, indexes[combination[0]-1], analyzer)
			firsts_b := FindProductionFirst(productions, indexes[combination[1]-1], analyzer)

			intersection := Intersection(firsts_a, firsts_b)

			// Primera regla: Interseccion de FIRST(a) y FIRST(b) debe ser un conjunto vacio.
			// la Segunda regla se verifica con esta misma condicion, ya que si existe mas de dos epsilon la interseccion
			// seria mayor a 0.
			if len(intersection) != 0 {
				return false
			}

			// Buscamos follows del no teminal
			follows := FindFollow(productions, non_terminal, analyzer, "")

			// Tercera regla.
			if Contains(firsts_a, "' '") {
				third_intersection := Intersection(firsts_b, follows)

				if len(third_intersection) != 0 {
					return false
				}
			}

			if Contains(firsts_b, "' '") {
				third_intersection := Intersection(firsts_a, follows)

				if len(third_intersection) != 0 {
					return false
				}
			}
		}
	}

	return true
}

func GetLL1Table(productions []string, analyzer AnalyzerOutput) {
	non_terminals := analyzer.non_terminals

	FIRSTS_CACHE_STATE = make(FIRSTS)

	// Initialize dictionary
	for _, non_terminal := range non_terminals {
		LL1_TABLE_STATE[non_terminal] = make(map[string]string)
	}

	for _, non_terminal := range non_terminals {

		for index, production := range productions {
			left, _ := SplitProduction(production)

			if left == non_terminal {
				firsts := FindProductionFirst(productions, index, analyzer)

				for _, first := range firsts {

					if first == "' '" {
						follows := FindFollow(productions, non_terminal, analyzer, "")
						for _, follow := range follows {
							new_production := fmt.Sprintf("%s -> ' '", non_terminal)
							LL1_TABLE_STATE[non_terminal][follow] = new_production
						}
					} else {
						LL1_TABLE_STATE[non_terminal][first] = production
					}
				}
			}
		}
	}

	// Initialize $
	analyzer.terminals = append(analyzer.terminals, "$")

	html_output := `
		<!DOCTYPE html>
		<html>
		<body>
		<table style="border: 1px solid black">
	`

	header_name := "Non Terminal"
	header := append([]string{header_name}, analyzer.terminals...)
	html_output += BuildHtmlRow(header, true)

	for key, values := range LL1_TABLE_STATE {
		var elements []string
		elements = append(elements, key)
		for index, terminal := range header {
			// Skip header title
			if index == 0 {
				continue
			}
			element := ""
			value, found := values[terminal]
			if found {
				element = value
			}
			elements = append(elements, element)
		}
		html_output += BuildHtmlRow(elements, false)
	}

	html_output += `
		</table>
		</body>
		</html>
		`

	output_file_dir := "output.html"
	file, err := os.Create(output_file_dir)

	if err != nil {
		fmt.Print("File could not be loaded")
		return
	} else {
		_, err := file.WriteString(html_output)

		if err != nil {
			fmt.Print("An error occured writing the file")
			return
		}
		fmt.Printf("\nFile saved succesfully: %s\n", output_file_dir)

	}
}

func CheckValidInput(productions []string, input string) bool {
	var stack []string

	left, _ := SplitProduction(productions[0])
	first_symbol := left

	stack = append(stack, "$")
	stack = append(stack, first_symbol)

	input_queue := strings.Split(input, " ")
	input_queue = append(input_queue, "$")

	for {
		// fmt.Printf("%s | %s\n", stack, input_queue)
		if len(stack) == 0 {
			return true
		}

		stack_element := LastElement(stack)
		queue_element := input_queue[0]

		if LastElement(stack) == queue_element {
			// POP
			stack = PopStack(stack)
			input_queue = PopQueue(input_queue)
			continue
		}

		rule, found := LL1_TABLE_STATE[stack_element][queue_element]
		if !found {
			return false
		}

		stack = PopStack(stack)

		_, right := SplitProduction(rule)

		if right == "' '" {
			continue
		}

		tokens := strings.Split(right, " ")
		tokens = Reverse(tokens)

		stack = append(stack, tokens...)
	}
}
