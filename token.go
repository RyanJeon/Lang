package main

import (
	"strings"
)

//Token represents a term in a line of code
type Token struct {
	Type  string
	Value []byte
}

func isInt(word string) bool {
	chars := strings.Split(word, "")

	//Loop over characters
	for _, c := range chars {
		bytechar := []byte(c)
		//one byte holds one character
		val := bytechar[0]

		//Not an integer 0:0x30 - 9:0x39
		if val < 48 || val > 57 {
			return false
		}
	}

	return true
}

func isOperator(word string) bool {
	chars := strings.Split(word, "")

	//Loop over characters
	for _, c := range chars {
		bytechar := []byte(c)
		//one byte holds one character
		val := bytechar[0]

		//Check if operator * : 42, + : 43, - : 45, / : 47
		if val == 44 || val == 46 || val < 40 || val > 47 {
			return false
		}
	}

	return true
}

func isDeclaration(word string) bool {
	return word == "Int"
}

func isAssignment(word string) bool {
	return word == "="
}

func tokenizer(line string) []Token {
	words := strings.Fields(line)
	tokens := []Token{}

	for i, w := range words {
		if isInt(w) {
			token := Token{
				Type:  "Int",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else if i > 0 && tokens[i-1].Type == "Declaration" {
			token := Token{
				Type:  "Variable",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else if isOperator(w) {
			token := Token{
				Type:  "Operator",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else if isDeclaration(w) {
			token := Token{
				Type:  "Declaration",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else if isAssignment(w) {
			token := Token{
				Type:  "Assignment",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else {
			_, ok := LocalVariable[w]
			//w is a local variable
			if ok {
				token := Token{
					Type:  "Variable",
					Value: []byte(w),
				}
				tokens = append(tokens, token)
			} else {
				token := Token{
					Type:  "Function",
					Value: []byte(w),
				}
				tokens = append(tokens, token)
			}
		}

	}

	return tokens
}
