package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

func asm64(tree *Tree) {
	f, _ := os.Create("test.asm")

	//Set up assembly
	f.Write([]byte(".data\n"))
	f.Write([]byte("print: .asciz \"\" \n"))
	f.Write([]byte(".globl execute\n")) //make main visible to linker
	f.Write([]byte(".text\n"))
	f.Write([]byte("execute:\n")) //main code segment
	treeAssemble(tree, f)
	f.Write([]byte("retq\n")) //end the process

}

func treeAssemble(tree *Tree, f *os.File) {
	//Operator means produce integer value ex) + - * ... etc
	if (*tree).Type == "Operator" {
		val := []byte(strconv.Itoa(Arithmetic(tree)))
		val = reverseByteArray(val)
		hex := hex.EncodeToString(val)
		//Move arithmetic output to rax
		s := fmt.Sprintf("movq	$0x%s, %%rcx\n", hex)
		f.Write([]byte(s))
	} else if (*tree).Type == "Function" {
		if string((*tree).Value) == "print" {
			treeAssemble(tree.Right, f)
			f.Write([]byte("movq	$0x2000004, %rax\n"))
			f.Write([]byte("movq	$1, %rdi\n"))
			f.Write([]byte("movq	%rcx, print(%rip)\n"))
			f.Write([]byte("leaq	print(%rip), %rsi\n"))
			f.Write([]byte("movq	$100, %rdx\n"))
			f.Write([]byte("syscall\n"))
		}
	}
}
