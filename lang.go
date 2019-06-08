package main

import "fmt"

func main() {
	example := "13912390123 + 1 + 20 - 40"
	tokenized := tokenizer(example)
	post := postfix(tokenized)
	fmt.Println(post)
}
