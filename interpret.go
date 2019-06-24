package main

import (
	"fmt"
	"log"
	"os"
)

//LocalVariable : map of local variable relative stack index
var LocalVariable map[string]int

//FunctionParamMap : Keeps track of how many parameter a function takes
var FunctionParamMap map[string]int

var stackindex int
var paramCount int //Keeps track of how many parameters are in a function

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
			offset, exist := LocalVariable[string((*tree).Value)]

			if !exist {
				err := fmt.Sprintf("Variable %s is not declared!", (*tree).Value)
				log.Fatal(err)
			}

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
	if tree == nil {
		return
	}
	termtype := (*tree).Type
	switch termtype {
	case "Declaration":
		Declaration(tree.Right, f, vartype)
		break
	//Should check if it's function call or declaration
	case "Function":
		code := fmt.Sprintf("%s:\n", string((*tree).Value))
		f.WriteString(code)
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
		f.WriteString("popq	%rbp\n")
		f.WriteString("pushq	%rax\n")
		break
	case "Operator":
		Arithmetic(tree, f)
		//move rbp to top of the stack again
		f.WriteString("popq	%rbp\n")
		f.WriteString("pushq	%rax\n")
		break
	}
}

//FunctionDeclaration : Writes assembly for function declaration statement
func FunctionDeclaration(tokens []Token, f *os.File) {
	code := fmt.Sprintf("%s:\n", string(tokens[1].Value))
	f.WriteString(code)
	f.WriteString("pushq   %rbp\n")
	f.WriteString("movq	%rsp, %rbp\n")
	// Dec Func ( **variables** ) {
	i := 3
	for i < len(tokens) {
		if tokens[i].Type == "Comma" {
			i++
			continue
		}
		if tokens[i].Type != "Declaration" || tokens[i+1].Type != "Variable" {
			log.Fatalf("Unexpected %s: Interpret error", string(tokens[i].Value))
		}
		LocalVariable[string(tokens[i+1].Value)] = stackindex - 8
		stackindex = stackindex + 8
		paramCount++ //Increase number of parameter
		i = i + 2
	}

	//Register # of parameters for the function
	FunctionParamMap[string(tokens[1].Value)] = paramCount
}

//FunctionCall : Provides assembly for functioncall statement
func FunctionCall(tokens []Token, f *os.File) {
	//Print exception need fix later
	if string(tokens[0].Value) == "print" {
		t := tree(TokensPostfix(tokens))
		asm64(&t, f)
		return
	}

	i := 2
	j := 2
	params := 0
	for i < len(tokens) && string(tokens[i].Value) != ")" {
		for j < len(tokens) && tokens[j].Type != "Comma" && string(tokens[j].Value) != ")" {
			j++
		}
		t := tree(TokensPostfix(tokens[i:j]))
		asm64(&t, f)
		j++
		i = j

		f.WriteString("pushq	%rax\n")
		params++
	}

	f.WriteString(fmt.Sprintf("callq	%s\n", string(tokens[0].Value)))
	f.WriteString(fmt.Sprintf("addq	$%d, %%rsp\n", params*8))
}
