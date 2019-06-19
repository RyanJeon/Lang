package main

import (
	"fmt"
	"os"
)

//LocalVariable : map of local variable relative stack index
var LocalVariable map[string]int

var stackindex int

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
		if (*tree).Type == "Variable" {
			//Look up stack index for the variable
			offset := LocalVariable[string((*tree).Value)]
			index := (len(LocalVariable)+1)*8 - offset
			code := fmt.Sprintf("movq	%d(%%rbp), %%rax\n", index)
			f.WriteString(code)
		} else {
			s := fmt.Sprintf("movq	$%s, %%rax\n", string((*tree).Value))
			f.WriteString(s)
		}
		//Plus
	} else if val == 43 {
		Arithmetic(tree.Left, f)
		f.WriteString("pushq	%rax\n")
		Arithmetic(tree.Right, f)
		f.WriteString("popq	%rcx\n")
		f.WriteString("addq	%rcx, %rax\n")
		//Minus
	} else if val == 45 {
		Arithmetic(tree.Left, f)
		f.WriteString("pushq	%rax\n")
		Arithmetic(tree.Right, f)
		f.WriteString("popq	%rcx\n")
		f.WriteString("subq	%rax, %rcx\n")
		f.WriteString("movq	%rcx, %rax\n")
		//Mult
	} else if val == 42 {
		Arithmetic(tree.Left, f)
		f.WriteString("pushq	%rax\n")
		Arithmetic(tree.Right, f)
		f.WriteString("popq	%rcx\n")
		f.WriteString("mulq	%rcx\n")
		//Div
	} else if val == 47 {
		Arithmetic(tree.Left, f)
		f.WriteString("pushq	%rax\n")
		Arithmetic(tree.Right, f)
		f.WriteString("movq	%rax, %rcx\n")
		f.WriteString("popq	%rax\n")
		f.WriteString("xor		%rdx, %rdx\n")
		f.WriteString("divq	%rcx\n")
	} else {

	}
}

//Declaration : traverse tree for variable declaration
func Declaration(tree *Tree, f *os.File, vartype string) {
	termtype := (*tree).Type
	switch termtype {
	case "Declaration":
		Declaration(tree.Right, f, vartype)
		break
	case "Assignment":
		Declaration(tree.Left, f, vartype)
		Declaration(tree.Right, f, vartype)

		f.WriteString("pushq	%rbp\n")
		f.WriteString("movq	%rsp, %rbp\n")
		break
	case "Variable":
		LocalVariable[string((*tree).Value)] = stackindex
		stackindex = stackindex + 8
		break
	case "Int":
		Arithmetic(tree, f)
		//move rbp to top of the stack again
		if len(LocalVariable) > 1 {
			//pop
			f.WriteString("popq	%rbp\n")
		}
		f.WriteString("pushq	%rax\n")
	}

}
