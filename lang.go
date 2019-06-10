package main

import "log"

func main() {
	example := "print ( ( 20 + 30 ) * 10 )"
	tokenized := tokenizer(example)
	post := postfix(tokenized)
	log.Print("PostFix: ")
	log.Println(post)

	t := tree(post)
	inorder(&t)
	asm64(&t)
}
