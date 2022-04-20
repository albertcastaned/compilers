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
