package main

import (
	"fmt"
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
		} else {
			code := fmt.Sprintf("callq	%s\n", string((*tree).Value))
			f.WriteString(code)
		}
		break
	//End of func but need to be scalable
	case "End":
		//If there was local variable
		if stackindex > 8 {

			f.WriteString("movq	%rbp, %rsp\n")
			f.WriteString("popq	%rbp\n") //restore rbp
			for i := 0; i < (stackindex/8)-1; i++ {
				//pop remaining local variable
				f.WriteString("popq	%rcx\n")
			}
			//Reset local variable map
			LocalVariable = make(map[string]int)
			stackindex = 8
		}
		f.WriteString("retq\n")
	}
}

//CodeGen takes a statement classification and outputs corresponding assembly
func CodeGen(class string, tokens []Token, f *os.File) {
	switch class {
	case "Test":
		t := tree(TokensPostfix(tokens))
		asm64(&t, f)
		break
	case "VariableDeclaration":
		t := tree(TokensPostfix(tokens))
		asm64(&t, f)
		break
	case "FunctionDeclaration":
		FunctionDeclaration(tokens, f)
		break
	case "FunctionCall":
		FunctionCall(tokens, f)
		break
	case "EndOfFunction":
		f.WriteString("movq	%rbp, %rsp\n")
		f.WriteString("popq	%rbp\n")

		//move rsp to point to the return address. (paramCount*8) is there to
		//take account of the fact that variables passed in as parameters are
		//above ret address in the stack, and local variables are right below
		//the return address. However, both types of variables are in LocalVariable
		//map
		code := fmt.Sprintf("addq	$%d, %%rsp\n", stackindex-8-(paramCount*8))
		f.WriteString(code)
		f.WriteString("retq\n")

		LocalVariable = make(map[string]int)
		stackindex = 8
		paramCount = 0 //reset param count for new function!
		break
	}
}
