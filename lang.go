package main

import "log"

func main() {
	example := "2 + 1 + 3 - 4"
	tokenized := tokenizer(example)
	post := postfix(tokenized)
	log.Println(post)
	t := tree(post)
	inorder(&t)
}
