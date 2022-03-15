package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const PRODUCTION_DELIMITER = " -> "

type AnalyzerOutput struct {
	terminals, non_terminals []string
}

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

func Contains(list []string, value string) bool {
	for _, temp := range list {
		if temp == value {
			return true
		}
	}
	return false
}

func GetTerminals(left []string, right []string) []string {
	var terminals []string

	right_size := len(right)
	for i := 0; i < right_size; i++ {
		tokens := strings.Split(right[i], " ")
		tokens_size := len(tokens)
		for j := 0; j < tokens_size; j++ {
			token := tokens[j]
			if !Contains(left, token) && token != "'" && token != "" {
				terminals = append(terminals, token)
			}
		}
	}

	return RemoveDuplicates(terminals)
}

func GetNonTerminals(left []string) []string {
	return RemoveDuplicates(left)
}

func Analyze(lines []string) AnalyzerOutput {
	var left []string
	var right []string

	size := len(lines)

	for i := 0; i < size; i++ {
		splitted := strings.Split(lines[i], PRODUCTION_DELIMITER)
		left = append(left, splitted[0])
		right = append(right, splitted[1])
	}

	nonTerminals := GetNonTerminals(left)
	terminals := GetTerminals(left, right)

	return AnalyzerOutput{terminals, nonTerminals}
}

func main() {
	var productions_amount int
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Scanf("%d", &productions_amount)
	var lines []string

	for i := 0; i < productions_amount; i++ {
		scanner.Scan()
		lines = append(lines, scanner.Text())
	}

	output := Analyze(lines)
	fmt.Printf(
		"Terminal: %s\nNon terminal: %s\n",
		strings.Join(output.terminals, ", "),
		strings.Join(output.non_terminals, ", "),
	)
}
