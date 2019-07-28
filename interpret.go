package main

import (
	"fmt"
	"log"
	"os"
)

//Call is interface for a function call
type Call struct {
	Name   string
	Inputs [][]Token
}

//LocalVariable : map of local variable relative stack index
var LocalVariable map[string]int

//FunctionParamMap : Keeps track of how many parameter a function takes
var FunctionParamMap map[string]int

//FunctionCallStack : Keeps track of function calls
var FunctionCallStack CallStack

//EndStack helps determine what "}" will be ending. ex) conditional, function, loop
var EndStack StringStack

//Edit
//IfStack stores blocks of statements. Should be empty at the end of compilation if all blocks were closed properly
var IfStack StringStack

//Edit
//IfEndStack stores the type of block for each block. ex) block under function is typed "Function", block under if conditional it is typed "If"
var IfEndStack StringStack

//BlockCounter stores how many statement blocks there are. (forloop and if)
var BlockCounter int

//Stack index should be different each block so implement stack!
var stackindex int

//LocalVariableCountStack is a stack that holds number of local variables in a block
var LocalVariableCountStack IntStack

//VariableStack holds declared variables to assist deleting local variables in blocks
var VariableStack StringStack

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
func Arithmetic(tree *Tree) {
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
			ASMWrite(code)
		} else if (*tree).Type == "FunctionCall" {
			var call Call
			FunctionCallStack, call = FunctionCallStack.Poll()
			FunctionCall(call)
		} else {
			s := fmt.Sprintf("movq	$%s, %%rax\n", string((*tree).Value))
			ASMWrite(s)
		}
		//Plus
	} else if val == 43 {
		Arithmetic(tree.Left)
		ASMWrite("pushq	%rax\n")
		Arithmetic(tree.Right)
		ASMWrite("popq	%rcx\n")
		ASMWrite("addq	%rcx, %rax\n")
		//Minus
	} else if val == 45 {
		Arithmetic(tree.Left)
		ASMWrite("pushq	%rax\n")
		Arithmetic(tree.Right)
		ASMWrite("popq	%rcx\n")
		ASMWrite("subq	%rax, %rcx\n")
		ASMWrite("movq	%rcx, %rax\n")
		//Mult
	} else if val == 42 {
		Arithmetic(tree.Left)
		ASMWrite("pushq	%rax\n")
		Arithmetic(tree.Right)
		ASMWrite("popq	%rcx\n")
		ASMWrite("mulq	%rcx\n")
		//Div
	} else if val == 47 {
		Arithmetic(tree.Left)
		ASMWrite("pushq	%rax\n")
		Arithmetic(tree.Right)
		ASMWrite("movq	%rax, %rcx\n")
		ASMWrite("popq	%rax\n")
		ASMWrite("xor		%rdx, %rdx\n")
		ASMWrite("divq	%rcx\n")
	} else {

	}
}

//Declaration : traverse tree for variable declaration
func Declaration(tree *Tree, vartype string) {
	if tree == nil {
		return
	}
	termtype := (*tree).Type
	switch termtype {
	case "Declaration":
		Declaration(tree.Right, vartype)
		break
	//Should check if it's function call or declaration
	case "Function":
		code := fmt.Sprintf("%s:\n", string((*tree).Value))
		ASMWrite(code)
		Declaration(tree.Right, vartype)
		break
	case "Assignment":
		Declaration(tree.Left, vartype)
		Declaration(tree.Right, vartype)

		ASMWrite("pushq	%rbp\n")
		ASMWrite("movq	%rsp, %rbp\n")
		break
	case "Variable":
		LocalVariable[string((*tree).Value)] = stackindex
		stackindex = stackindex + 8
		break
	case "Int":
		Arithmetic(tree)
		//move rbp to top of the stack again
		ASMWrite("popq	%rbp\n")
		ASMWrite("pushq	%rax\n")
		break
	case "Operator":
		Arithmetic(tree)
		//move rbp to top of the stack again
		ASMWrite("popq	%rbp\n")
		ASMWrite("pushq	%rax\n")
		break
	}
}

//VariableDeclaration : a helper for declaring variable
func VariableDeclaration(tokens []Token) {

	variableName := string(tokens[1].Value)
	newTokenList := make([]Token, 0)
	//Currently just in type so just do arithmetic
	for i := 3; i < len(tokens); i++ {
		//If function call seen
		if tokens[i].Type == "Function" {
			j := i
			for j < len(tokens) && string(tokens[j].Value) != ")" {
				j++
			}
			//Make a function call
			AddFunctionCallToStack(tokens[i : j+1])

			tokens[i].Type = "FunctionCall"
			newTokenList = append(newTokenList, tokens[i])
			i = j
		} else {
			newTokenList = append(newTokenList, tokens[i])
		}
	}

	TokenProcess(newTokenList)

	//Increment variable counts
	if LocalVariableCountStack.isEmpty() {
		LocalVariableCountStack = LocalVariableCountStack.Push(1)
	} else {
		var count int
		//Increment number of local variables in current block
		LocalVariableCountStack, count = LocalVariableCountStack.Pop()
		LocalVariableCountStack = LocalVariableCountStack.Push(count + 1)
	}

	log.Println(LocalVariableCountStack)
	//Push declared variable onto the variablestack
	VariableStack = VariableStack.Push(variableName)

	LocalVariable[variableName] = stackindex
	stackindex = stackindex + 8

	ASMWrite("popq	%rbp\n")
	ASMWrite("pushq	%rax\n")

	ASMWrite("pushq	%rbp\n")
	ASMWrite("movq	%rsp, %rbp\n")
}

//FunctionDeclaration : Writes assembly for function declaration statement
func FunctionDeclaration(tokens []Token) {
	code := fmt.Sprintf("%s:\n", string(tokens[1].Value))
	ASMWrite(code)
	ASMWrite("pushq   %rbp\n")
	ASMWrite("movq	%rsp, %rbp\n")
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

	blockInit()
	//Register # of parameters for the function
	FunctionParamMap[string(tokens[1].Value)] = paramCount
}

//FunctionCall : Provides assembly for functioncall statement
func FunctionCall(functionCall Call) {
	function := functionCall.Name
	inputs := functionCall.Inputs

	params := 0
	for _, input := range inputs {
		TokenProcess(input)

		//Push the function argument to rax after processing current arg
		ASMWrite("pushq	%rax\n")
		params++
	}

	ASMWrite(fmt.Sprintf("callq	%s\n", function))
	ASMWrite(fmt.Sprintf("addq	$%d, %%rsp\n", params*8))
}

//AddFunctionCallToStack : adds function call to call stack
func AddFunctionCallToStack(tokens []Token) Call {
	callInputs := make([][]Token, 0)

	i := 2
	j := 2

	//Iterate through the input for the function call
	for i < len(tokens) && string(tokens[i].Value) != ")" {
		for j < len(tokens) && tokens[j].Type != "Comma" && string(tokens[j].Value) != ")" {
			j++
		}

		//One input in the function call
		input := tokens[i:j]
		callInputs = append(callInputs, input)

		j++
		i = j
	}

	call := Call{
		Name:   string(tokens[0].Value),
		Inputs: callInputs,
	}
	FunctionCallStack = FunctionCallStack.Push(call)
	return call
}

//FunctionReturn is used to return value for function
func FunctionReturn(tokens []Token) {

	newTokenList := make([]Token, 0)
	//Currently just in type so just do arithmetic
	for i := 0; i < len(tokens); i++ {
		//If function call seen
		if tokens[i].Type == "Function" {
			j := i
			for j < len(tokens) && string(tokens[j].Value) != ")" {
				j++
			}
			//Make a function call
			AddFunctionCallToStack(tokens[i : j+1])

			tokens[i].Type = "FunctionCall"
			newTokenList = append(newTokenList, tokens[i])
			i = j
		} else {
			newTokenList = append(newTokenList, tokens[i])
		}
	}

	TokenProcess(newTokenList[1:])

	//Pop local variable, restore rsp to return address
	FunctionEndRspReset()
}

func isCondOp(token Token) bool {
	word := string(token.Value)
	return word == "or" || word == "and" || word == ">" || word == "<" || word == "==" || word == "!="
}

//takes conditional statement and generates assembly
func conditionalHelper(tokens []Token) {

	//Left and right pointer to parse portion of conditional statement
	left := 0
	right := 0

	//general conditional statement grammar :
	// {expression} {conditional operator} {expression}
	// Need to check that expressions on both sides are the same type

	for i, token := range tokens {

		//If current token is conditional operator
		if isCondOp(token) {
			right = i + 1
			for right < len(tokens) && !isCondOp(tokens[right]) {
				right++
			}

			//segment of conditional expression
			lhs := tokens[left:i]
			rhs := tokens[i+1 : right]
			op := tokens[i]

			conditionalExpGen(lhs, rhs, op)
			left = right
		}
	}
}

func conditionalExpGen(lhs []Token, rhs []Token, op Token) {

	TokenProcess(lhs)
	ASMWrite("movq	%rax, %rbx\n")
	TokenProcess(rhs)
	ASMWrite("cmpq	%rax, %rbx\n")

	//If the condition in if statement is not met, it will jump to the "jump" block
	jump := fmt.Sprintf("ifblock_%d", BlockCounter)

	//do everything within the if block
	// Here
	////////

	switch string(op.Value) {
	case ">":
		code := fmt.Sprintf("jle	%s\n", jump)
		ASMWrite(code)
	case "<":
		code := fmt.Sprintf("jge	%s\n", jump)
		ASMWrite(code)
	case "==":
		code := fmt.Sprintf("jne	%s\n", jump)
		ASMWrite(code)
	case "!=":
		code := fmt.Sprintf("je	%s\n", jump)
		ASMWrite(code)
	}

	//Push the jump address to the if stack
	IfStack = IfStack.Push(jump)

}

//IfStatement : When the keyword IF is detected
func IfStatement(tokens []Token) {
	if string(tokens[len(tokens)-1].Value) != "=>" {
		log.Fatal("Expected => in if statement")
	} else {
		conditional := tokens[1 : len(tokens)-1]
		conditionalHelper(conditional)
	}

	blockInit()

	//Increment block counter to avoid conflict
	BlockCounter++
}

//ElseStatement : When you encounter else statement
func ElseStatement(f *os.File) {
	var address string
	IfStack, address = IfStack.Pop()

	code := fmt.Sprintf("%s:\n", address)
	ASMWrite(code)
}

//IfEnd for ending an if conditional
func IfEnd() {
	var address string
	IfStack, address = IfStack.Pop()

	code := fmt.Sprintf("%s:\n", address)
	ASMWrite(code)

	//If block has been executed jump to the end of the if statement
	// ifEnd := fmt.Sprintf("ifEnd_%d", BlockCounter)
	// code = fmt.Sprintf("jmp	%s\n", ifEnd)
	// ASMWrite(code)
	// IfEndStack = IfEndStack.Push(ifEnd)
	// BlockCounter++
}

//RedefineVariable to redefine already declared variable. Note: Could be used for declaring too
func RedefineVariable(tokens []Token) {
	variableName := string(tokens[0].Value)
	newTokenList := make([]Token, 0)
	//Currently just in type so just do arithmetic
	for i := 2; i < len(tokens); i++ {
		//If function call seen
		if tokens[i].Type == "Function" {
			j := i
			for j < len(tokens) && string(tokens[j].Value) != ")" {
				j++
			}
			//Make a function call
			AddFunctionCallToStack(tokens[i : j+1])

			tokens[i].Type = "FunctionCall"
			newTokenList = append(newTokenList, tokens[i])
			i = j
		} else {
			newTokenList = append(newTokenList, tokens[i])
		}
	}

	TokenProcess(newTokenList)

	offset := LocalVariable[variableName]

	index := (len(LocalVariable)+1)*8 - offset
	code := fmt.Sprintf("movq	%%rax, %d(%%rbp)\n", index)
	ASMWrite(code)
}

//Initialize local variable count for current block
func blockInit() {
	LocalVariableCountStack = LocalVariableCountStack.Push(0)
}

//When block ends, all the local variables declared inside needs to be deleted
func BlockEnd() {
	var count int
	LocalVariableCountStack, count = LocalVariableCountStack.Pop()

	//Pop count amount of local variables
	PopLocalVariables(count)
	//If there was at least one local variable at the current block
	if count != 0 {
		log.Println(LocalVariable)
		for count != 0 {
			var variable string
			VariableStack, variable = VariableStack.Pop()
			stackindex = stackindex - 8

			//Delete current variable from the map to make sure the variable cannot be referenced again
			delete(LocalVariable, variable)
			count--
		}

		log.Println(LocalVariable)
	}
}
