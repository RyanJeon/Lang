package main

import (
	"bufio"
	"fmt"
)

//KEEP ALL THE RESULTS TO RAX!
//For arithmetic for now

//TokenProcess general token prcessing function
func TokenProcess(tokens []Token) {

	//Take tokens, covert them into RPN, and produce AST
	t := tree(TokensPostfix(tokens))

	//Get the syntax tree and output assembly
	treeAssemble(&t)
}

func treeAssemble(tree *Tree) {
	switch (*tree).Type {
	//Operator means produce integer value ex) + - * ... etc
	case "Operator":
		Arithmetic(tree)
		break
	case "Int":
		Arithmetic(tree)
		break
	case "Variable":
		//Should change to allow dif variable types
		Arithmetic(tree)
		break
	//Fix Later
	case "FunctionCall":
		//Should change to allow dif variable types
		Arithmetic(tree)
		break
	case "Declaration":
		Declaration(tree, string((*tree).Value))
		break
	//This is function call
	case "Function":
		if string((*tree).Value) == "출력" || string((*tree).Value) == "print" {
			treeAssemble(tree.Right)
			ASMWrite("callq	inttostring\n")
			ASMWrite("callq	printout\n")
		} else {
			code := fmt.Sprintf("callq	%s\n", string((*tree).Value))
			ASMWrite(code)
		}
		break
	//End of func but need to be scalable
	case "}":
		//If there was local variable
		if stackindex > 8 {

			ASMWrite("movq	%rbp, %rsp\n")
			ASMWrite("popq	%rbp\n") //restore rbp
			for i := 0; i < (stackindex/8)-1; i++ {
				//pop remaining local variable
				ASMWrite("popq	%rcx\n")
			}
			//Reset local variable map
			LocalVariable = make(map[string]int)
			stackindex = 8
		}
		ASMWrite("retq\n")
	}
}

//ScanAndGen takes a scanner object and generates code given string input
func ScanAndGen(scanner *bufio.Scanner) {
	for scanner.Scan() {
		code := scanner.Text()
		if len(code) == 0 {
			continue
		}
		tokenized := tokenizer(code)
		class := ClassifyStatement(tokenized)
		CodeGen(class, tokenized, scanner)
	}
}

//CodeGen takes a statement classification and outputs corresponding assembly
func CodeGen(class string, tokens []Token, scanner *bufio.Scanner) {
	switch class {
	case "Test":
		TokenProcess(tokens)
		break
	case "VariableDeclaration":
		VariableDeclaration(tokens)
		break
	case "FunctionReturn":
		FunctionReturn(tokens)
	case "FunctionDeclaration":
		EndStack = EndStack.Push("FunctionDeclaration")
		FunctionDeclaration(tokens)
		break
	case "FunctionCall":
		//Print exception need fix later
		if string(tokens[0].Value) == "print" {
			TokenProcess(tokens)
		} else {
			FunctionCall(AddFunctionCallToStack(tokens))
			FunctionCallStack, _ = FunctionCallStack.Pop()
		}
		break
	case "EndOf":
		var endType string
		//Check what is ending
		EndStack, endType = EndStack.Pop()

		if endType == "IfStatement" {
			BlockEnd()
			IfEnd()
		} else if endType == "FunctionDeclaration" {

			//Pop local variable, restore rsp to return address
			FunctionEndRspReset()

			LocalVariable = make(map[string]int)
			FunctionCallStack = make([]Call, 0)
			stackindex = 8
			paramCount = 0 //reset param count for new function!
		}

		break
	case "IfStatement":
		EndStack = EndStack.Push("IfStatement")
		IfStatement(tokens)
	case "Redefinition":
		RedefineVariable(tokens)
	}
}
