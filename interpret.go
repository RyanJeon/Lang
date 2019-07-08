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

//endStack helps determine what "}" will be ending. ex) conditional, function, loop
var EndStack StringStack

//Edit
//IfStack stores blocks of statements. Should be empty at the end of compilation if all blocks were closed properly
var IfStack StringStack

//Edit
//IfEndStack stores the type of block for each block. ex) block under function is typed "Function", block under if conditional it is typed "If"
var IfEndStack StringStack

//BlockCounter stores how many statement blocks there are. (forloop and if)
var BlockCounter int

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
		} else if (*tree).Type == "FunctionCall" {
			var call Call
			FunctionCallStack, call = FunctionCallStack.Poll()
			FunctionCall(call, f)
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

func VariableDeclaration(tokens []Token, f *os.File) {
	variableName := string(tokens[1].Value)
	LocalVariable[variableName] = stackindex
	stackindex = stackindex + 8

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

	t := tree(TokensPostfix(newTokenList))
	Arithmetic(&t, f)

	f.WriteString("popq	%rbp\n")
	f.WriteString("pushq	%rax\n")

	f.WriteString("pushq	%rbp\n")
	f.WriteString("movq	%rsp, %rbp\n")
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
func FunctionCall(functionCall Call, f *os.File) {
	function := functionCall.Name
	inputs := functionCall.Inputs

	params := 0
	for _, input := range inputs {
		t := tree(TokensPostfix(input))
		asm64(&t, f)
		f.WriteString("pushq	%rax\n")
		params++
	}

	f.WriteString(fmt.Sprintf("callq	%s\n", function))
	f.WriteString(fmt.Sprintf("addq	$%d, %%rsp\n", params*8))
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

func FunctionReturn(tokens []Token, f *os.File) {
	t := tree(TokensPostfix(tokens[1:]))
	asm64(&t, f)
}

//IfStatement : When the keyword IF is detected
func IfStatement(tokens []Token, f *os.File) {
	f.WriteString("cmpq	%rax, %rbx\n")

	//If the condition in if statement is not met, it will jump to the "jump" block
	jump := fmt.Sprintf("ifblock_%d", BlockCounter)
	code := fmt.Sprintf("jmp	%s\n", jump)
	f.WriteString(code)
	//Push the jump address to the if stack
	IfStack = IfStack.Push(jump)

	//do everything within the if block
	// Here
	////////

	//Increment block counter to avoid conflict
	BlockCounter++
}

//IfEnd for ending an if conditional
func IfEnd(f *os.File) {
	var address string
	IfStack, address = IfStack.Pop()

	code := fmt.Sprintf("%s:\n", address)
	f.WriteString(code)
	//If block has been executed jump to the end of the if statement
	// ifEnd := fmt.Sprintf("ifEnd_%d", BlockCounter)
	// code = fmt.Sprintf("jmp	%s\n", ifEnd)
	// f.WriteString(code)
	// IfEndStack = IfEndStack.Push(ifEnd)
	// BlockCounter++
}
