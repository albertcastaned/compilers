package main

import "strings"

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
