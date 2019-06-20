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
	case "Function":
		if string((*tree).Value) == "출력" || string((*tree).Value) == "print" {
			treeAssemble(tree.Right, f)
			f.WriteString("callq	inttostring\n")
			f.WriteString("callq	printout\n")
		}
		break
	}
}
