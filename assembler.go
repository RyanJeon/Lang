package main

import (
	"os"
)

func asm64(tree *Tree, f *os.File) {
	//Set up assembly
	treeAssemble(tree, f)
}

func treeAssemble(tree *Tree, f *os.File) {

	switch (*tree).Type {
	//Operator means produce integer value ex) + - * ... etc
	case "Operator":
		Arithmetic(tree, f)
		break
	case "Int":
		Arithmetic(tree, f)
		break
	case "Variable":
		//Should change to allow dif variable types
		Arithmetic(tree, f)
		break
	case "Declaration":
		Declaration(tree, f, string((*tree).Value))
		break
	//This is function call
	case "Function":
		if string((*tree).Value) == "출력" || string((*tree).Value) == "print" {
			treeAssemble(tree.Right, f)
			f.WriteString("callq	inttostring\n")
			f.WriteString("callq	printout\n")
		}
		break
	//End of func but need to be scalable
	case "CurlyLeft":
		f.WriteString("movq	%rbp, %rsp\n")
		f.WriteString("popq	%rbp\n") //restore rbp
		for i := 0; i < (stackindex/8)-1; i++ {
			//pop remaining local variable
			f.WriteString("popq	%rcx\n")
		}
		f.WriteString("retq\n")
		//Reset local variable map
		LocalVariable = make(map[string]int)
		stackindex = 8
	}
}
