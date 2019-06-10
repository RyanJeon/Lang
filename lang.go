package main

import (
	"log"
)

func main() {
	example := "( 3 + 10 ) * ( 2 - 3 ) / 5"
	tokenized := tokenizer(example)
	post := postfix(tokenized)
	log.Println(post)
	t := tree(post)
	log.Println(interpret(&t))
	log.Println(interpret(&t))
	asm64(&t)
}
