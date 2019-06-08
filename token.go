package main

import (
	"strings"
)

type Token struct {
	Type  string
	Value []byte
}

func isInt(word string) bool {
	chars := strings.Split(word, "")

	//Loop over characters
	for _, c := range chars {
		byte_char := []byte(c)
		//one byte holds one character
		val := byte_char[0]

		//Not an integer 0:0x30 - 9:0x39
		if val < 48 || val > 57 {
			return false
		}
	}

	return true
}

func tokenizer(line string) []Token {
	words := strings.Fields(line)
	tokens := []Token{}

	for _, w := range words {
		if isInt(w) {
			token := Token{
				Type:  "Int",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		} else {
			token := Token{
				Type:  "Operator",
				Value: []byte(w),
			}
			tokens = append(tokens, token)
		}
	}

	return tokens
}
