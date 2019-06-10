package main

import (
	"fmt"
	"os"
	"strconv"
)

func asm64(tree *Tree) {
	f, _ := os.Create("test.asm")

	//Set up assembly
	f.Write([]byte(".data\n"))
	f.Write([]byte(".globl execute\n")) //make main visible to linker
	f.Write([]byte(".text\n"))
	f.Write([]byte("execute:\n")) //main code segment
	treeAssemble(tree, f, "")
	f.Write([]byte("retq\n")) //end the process

}

func treeAssemble(tree *Tree, f *os.File, prevType string) {
	fmt.Println((*tree).Type)
	fmt.Println(prevType)
	if (*tree).Type != "Operator" {
		s := ""
		//Second int in arithmetic exp
		if prevType == "rbx" {
			s = fmt.Sprintf("movq	$%s, %%rbx\n", (*tree).Value)
		} else {
			s = fmt.Sprintf("movq	$%s, %%rax\n", (*tree).Value)
		}
		f.Write([]byte(s))
	} else { //Operator means produce integer value ex) + - * ... etc
		val := strconv.Itoa(interpret(tree))
		s := fmt.Sprintf("movq	$%s, %%rax\n", []byte(val))
		f.Write([]byte(s))
	}
}
