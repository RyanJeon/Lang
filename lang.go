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

	//Initialize Local Variable map
	LocalVariable = make(map[string]int)
	//Initialize stack index for local variables
	stackindex = 8

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, _ := os.Create("test.asm")
	f.Write([]byte(".data\n"))
	f.Write([]byte("print: .asciz \"\" \n"))
	f.Write([]byte(".globl execute\n")) //make main visible to linker
	f.Write([]byte(".text\n"))
	IntToHex(f)
	PrintAsm(f)
	f.Write([]byte("execute:\n")) //main code segment

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		tokenized := tokenizer(code)
		post := postfix(tokenized)
		t := tree(post)
		asm64(&t, f)
	}

	f.Write([]byte("movq	$42, %rbx\n")) //end the process
	f.Write([]byte("syscall\n"))
}
