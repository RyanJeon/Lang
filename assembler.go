package main

import (
	"os"
)

func asm64(tree *Tree, f *os.File) {
	//Set up assembly
	treeAssemble(tree, f)
}

func treeAssemble(tree *Tree, f *os.File) {
	//Operator means produce integer value ex) + - * ... etc
	if (*tree).Type == "Operator" {
		//Do arithmetic if operator is seen
		Arithmetic(tree, f)
	} else if (*tree).Type == "Function" {
		if string((*tree).Value) == "print" {
			treeAssemble(tree.Right, f)
			f.WriteString("callq	inttohex\n")
			f.WriteString("callq	printout\n")
		}
	}
}
