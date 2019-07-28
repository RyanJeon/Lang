package main

import (
	"bufio"
	"log"
	"os"
)

//File created to hold assembly
var File *os.File

func main() {
	// example := "print ( ( 20 + 30 ) * 10 + ( 2 - 65 * 10 ) * 32 )"
	args := os.Args
	file, err := os.Open(args[1])

	//Initialize Local Variable map
	LocalVariable = make(map[string]int)
	FunctionParamMap = make(map[string]int)
	//Initialize stack index for local variables
	stackindex = 8

	GrammarInit()

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	File, _ = os.Create("test.asm")
	ASMInit()

	scanner := bufio.NewScanner(file)
	ScanAndGen(scanner)
}
