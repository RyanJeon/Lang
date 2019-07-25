package main

import (
	"bufio"
	"fmt"
	"os"
)

//TokenProcess general token prcessing function
func TokenProcess(tokens []Token, f *os.File) {
	t := tree(TokensPostfix(tokens))
	asm64(&t, f)
}

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
	//Fix Later
	case "FunctionCall":
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
	case "}":
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

//ScanAndGen takes a scanner object and generates code given string input
func ScanAndGen(scanner *bufio.Scanner, f *os.File) {
	for scanner.Scan() {
		code := scanner.Text()
		if len(code) == 0 {
			continue
		}
		tokenized := tokenizer(code)
		class := ClassifyStatement(tokenized)
		CodeGen(class, tokenized, f, scanner)
	}
}

//CodeGen takes a statement classification and outputs corresponding assembly
func CodeGen(class string, tokens []Token, f *os.File, scanner *bufio.Scanner) {
	switch class {
	case "Test":
		t := tree(TokensPostfix(tokens))
		asm64(&t, f)
		break
	case "VariableDeclaration":
		VariableDeclaration(tokens, f)
		break
	case "FunctionReturn":
		FunctionReturn(tokens, f)
	case "FunctionDeclaration":
		EndStack = EndStack.Push("FunctionDeclaration")
		FunctionDeclaration(tokens, f)
		break
	case "FunctionCall":
		//Print exception need fix later
		if string(tokens[0].Value) == "print" {
			t := tree(TokensPostfix(tokens))
			asm64(&t, f)
		} else {
			FunctionCall(AddFunctionCallToStack(tokens), f)
			FunctionCallStack, _ = FunctionCallStack.Pop()
		}
		break
	case "EndOf":
		var endType string
		//Check what is ending
		EndStack, endType = EndStack.Pop()

		if endType == "IfStatement" {
			IfEnd(f)
		} else if endType == "FunctionDeclaration" {
			f.WriteString("movq	%rbp, %rsp\n")
			f.WriteString("popq	%rbp\n")

			//move rsp to point to the return address. (-paramCount) is there to
			//take account of the fact that variables passed in as parameters are
			//above ret address in the stack, and local variables are right below
			//the return address. However, both types of variables are in LocalVariable
			//map meaning, len(LocalVariable) will count both types of variables!
			//
			// [         ]
			// [  param  ]
			// [         ]
			// [ ret ad  ]
			// [         ]
			// [local var]
			// [         ]  <== rsp

			//Note : how do we deal with resetting rsp when there is variable declaration
			//that is not hit?
			code := fmt.Sprintf("addq	$%d, %%rsp\n", (len(LocalVariable)-paramCount)*8)
			f.WriteString(code)
			f.WriteString("retq\n")

			LocalVariable = make(map[string]int)
			FunctionCallStack = make([]Call, 0)
			stackindex = 8
			paramCount = 0 //reset param count for new function!
		}

		break
	case "IfStatement":
		EndStack = EndStack.Push("IfStatement")
		IfStatement(tokens, f)
	case "Redefinition":
		RedefineVariable(tokens, f)
	}
}
