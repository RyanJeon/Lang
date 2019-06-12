package main

import (
	"fmt"
	"os"
)

//BytesToInt : Convert int ascii value to int
func BytesToInt(bytes []byte) int {
	//0 is 0x30
	zero := 48
	res := 0
	for _, c := range bytes {
		res = res*10 + (int(c) - zero)
	}
	return res
}

//NOTE: Fix Negatives
//Arithmetic : given Arithmetic tree, calculates integer value
func Arithmetic(tree *Tree, f *os.File) {
	val := (*tree).Value[0]
	if (*tree).Type != "Operator" {
		s := fmt.Sprintf("movq	$%s, %%rax\n", string((*tree).Value))
		f.Write([]byte(s))
		//Plus
	} else if val == 43 {
		Arithmetic(tree.Left, f)
		f.Write([]byte("pushq	%rax\n"))
		Arithmetic(tree.Right, f)
		f.Write([]byte("popq	%rcx\n"))
		f.Write([]byte("addq	%rcx, %rax\n"))
		//Minus
	} else if val == 45 {
		Arithmetic(tree.Left, f)
		f.Write([]byte("pushq	%rax\n"))
		Arithmetic(tree.Right, f)
		f.Write([]byte("popq	%rcx\n"))
		f.Write([]byte("subq	%rax, %rcx\n"))
		f.Write([]byte("movq	%rcx, %rax\n"))
		//Mult
	} else if val == 42 {
		Arithmetic(tree.Left, f)
		f.Write([]byte("pushq	%rax\n"))
		Arithmetic(tree.Right, f)
		f.Write([]byte("popq	%rcx\n"))
		f.Write([]byte("mulq	%rcx\n"))
		//Div
	} else if val == 47 {
		Arithmetic(tree.Left, f)
		f.Write([]byte("pushq	%rax\n"))
		Arithmetic(tree.Right, f)
		f.Write([]byte("movq	%rax, %rcx\n"))
		f.Write([]byte("popq	%rax\n"))
		f.Write([]byte("xor		%rdx, %rdx\n"))
		f.Write([]byte("divq	%rcx\n"))
	} else {

	}
}
