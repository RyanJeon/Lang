package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	// example := "print ( ( 20 + 30 ) * 10 + ( 2 - 65 * 10 ) * 32 )"
	args := os.Args
	file, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		tokenized := tokenizer(code)
		post := postfix(tokenized)
		t := tree(post)
		asm64(&t)
	}
}
