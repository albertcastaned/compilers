package main

import (
	"fmt"
	"strings"
)

// Funcion para regresar una lista de cadenas libre de duplicados dado
// una lista de cadenas.
func RemoveDuplicates(str []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range str {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Funcion que dado una lista de cadenas y valor, se
// busque si existe el valor en la lista especificada
func Contains(list []string, value string) bool {
	for _, temp := range list {
		if temp == value {
			return true
		}
	}
	return false
}

// Funcion que dado una produccion ,regresa su izquierdo y dercho
func SplitProduction(production string) (string, string) {
	splitted := strings.Split(production, PRODUCTION_DELIMITER)
	left := splitted[0]
	right := splitted[1]

	return left, right
}

// Funcion para eliminar los epsilons de una lista de cadendas
func RemoveEpsilons(s []string) []string {
	for i, v := range s {
		if v == "' '" {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Definimos el tipo de dato Combination de dos valores enteros.
type Combination [2]int

// Funcion para crear un arreglo de combinaciones posibles numericos de 1 a n
func CreateCombinations(n int) []Combination {
	var result []Combination

	for left_index := 1; left_index < n+1; left_index++ {
		for right_index := left_index + 1; right_index < n+1; right_index++ {
			var combination Combination
			combination[0] = left_index
			combination[1] = right_index
			result = append(result, combination)
		}
	}
	return result
}

// Funcion para obtener la interseccion de dos arreglos.
func Intersection(a, b []string) (c []string) {
	result := make(map[string]bool)

	for _, item := range a {
		result[item] = true
	}

	for _, item := range b {
		if _, ok := result[item]; ok {
			c = append(c, item)
		}
	}
	return
}

// Funcion que convierte un arreglo de valores a un renglon de tabla de HTML
func BuildHtmlRow(values []string, isHeader bool) string {
	var rowStart, rowEnd, result string
	result = "<tr>"

	if isHeader {
		rowStart = "<th style=\"border: 1px solid black\">"
		rowEnd = "</th>"
	} else {
		rowStart = "<td style=\"border: 1px solid black\">"
		rowEnd = "</td>"
	}

	for _, value := range values {
		result += fmt.Sprintf("\n%s %s %s\n", rowStart, value, rowEnd)
	}
	result += "</tr>\n"

	return result
}

func PopStack(array []string) []string {
	return array[:len(array)-1]
}

func PopQueue(array []string) []string {
	_, new_array := array, array[1:]
	return new_array
}

func LastElement(array []string) string {
	return array[len(array)-1]
}

func Reverse(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return ss
}
