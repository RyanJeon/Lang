package main

import (
	"fmt"
	"os"
)

func asm64(tree *Tree) {
	f, _ := os.Create("compiled.asm")

	//Set up assembly
	f.Write([]byte(".globl _main\n")) //make main visible to linker
	f.Write([]byte("_main:\n"))       //main code segment
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
	} else {
		//Plus
		if (*tree).Value[0] == 43 {
			treeAssemble(tree.Left, f, "rax")
			treeAssemble(tree.Right, f, "rbx")
			f.Write([]byte("addq	%rbx, %rax\n"))
		} else if (*tree).Value[0] == 45 {
			treeAssemble(tree.Left, f, "rax")
			treeAssemble(tree.Right, f, "rbx")
			f.Write([]byte("subq	%rbx, %rax\n"))
		} else {

		}

	}
}
