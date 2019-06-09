package main

import (
	"log"
)

func main() {
	example := "2 + 4 + 60"
	tokenized := tokenizer(example)
	post := postfix(tokenized)
	log.Println(post)
	t := tree(post)
	// inorder(&t)
	log.Println(interpret(&t))
	asm64(&t)
}
