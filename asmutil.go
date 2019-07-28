package main

import (
	"fmt"
	"os"
)

//Assembly constants
const (
	MOVEQ = "movq"
	LEAQ  = "leaq"
	RAX   = "%rax"
	RDI   = "%rdi"
	RSI   = "%rsi"
	RIP   = "%rip"
)

//ToDo XOR used resistors to clean up just to be sure

//ASMInit : initialize assembly file
func ASMInit(f *os.File) {
	f.Write([]byte(".data\n"))
	f.Write([]byte("print: .asciz \"\" \n"))
	f.Write([]byte(".globl execute\n")) //make main visible to linker
	f.Write([]byte(".text\n"))
	intToHex(f)
	printAsm(f)
}

//IntToHex : writes integer to decimal equivalent of ascii in current file takes rax as input, and outputs to rcx
func intToHex(f *os.File) {
	//flip the digits before conversion
	f.WriteString("inttostring:\n")
	//Initialize variables, r9 = 10 (digit iterator), cl = 0, (8*n) bit shift
	f.WriteString("movq	$10, %r9\n")
	f.WriteString("xor	%rbx, %rbx\n")

	//When done, add end of line
	f.WriteString("movq	$0xa, %rbx\n")
	f.WriteString("salq	$8, %rbx\n")

	//start of adding digits
	f.WriteString("divide:\n")

	//Evaluate next digit, add digit * 16^n for equivalent ascii value
	f.WriteString("xor		%rdx, %rdx\n")
	f.WriteString("divq	%r9\n")         //val / 10
	f.WriteString("pushq	%rax\n")       //Save current int value
	f.WriteString("movq	$0x30, %rax\n") //Move zero char (0x30) to rax
	f.WriteString("addq	%rdx, %rax\n")  //add the remainder (next digit) to rax
	f.WriteString("salq	$8, %rbx\n")    // shift 8*n bits to left
	f.WriteString("addq	%rax, %rbx\n")  // add the current digit to rcx
	f.WriteString("popq	%rax\n")        //restore current val
	f.WriteString("test	%rax, %rax\n")  //Loop back if not done
	f.WriteString("jnz		divide\n")

	//Clear rcx, and move final ascii value to rcx
	f.WriteString("xor		%rcx, %rcx\n")
	f.WriteString("movq	%rbx, %rcx\n")

	f.WriteString("retq\n")

}

//Takes rax and flips its digits
func flipDigits(f *os.File) {
	f.WriteString("xor		%rcx, %rcx\n") //rcx will store result
	f.WriteString("movq	$10, %r9\n")
	f.WriteString("flipdigit:\n")
	f.WriteString("xor		%rdx, %rdx\n")
	f.WriteString("divq	%r9\n")
	f.WriteString("pushq	%rax\n")      //save remaining value
	f.WriteString("pushq	%rdx\n")      //save remainder
	f.WriteString("movq	%rcx, %rax\n") //Move return value to rax
	f.WriteString("mulq	%r9\n")        //return * 10
	f.WriteString("popq	%rdx\n")       //restore remainder
	f.WriteString("addq	%rdx, %rax\n") //add remainder to the result
	f.WriteString("movq	%rax, %rcx\n") //save the result back to rcx
	f.WriteString("popq	%rax\n")       //restore remaining value
	f.WriteString("test	%rax, %rax\n") //Loop back if not done
	f.WriteString("jnz		flipdigit\n")
	f.WriteString("movq	%rcx, %rax\n") //move the result to rax
}

func printAsm(f *os.File) {
	f.WriteString("printout:\n")
	f.WriteString("movq	$0x2000004, %rax\n")
	f.WriteString("movq	$1, %rdi\n")
	f.WriteString("movq	%rcx, print(%rip)\n")
	f.WriteString("leaq	print(%rip), %rsi\n")
	f.WriteString("movq	$100, %rdx\n")
	f.WriteString("syscall\n")
	f.WriteString("retq\n")
}

//FunctionEndRspReset : pop local varibles through resetting rsp based on length of local variable map
func FunctionEndRspReset() {
	File.WriteString("movq	%rbp, %rsp\n")
	File.WriteString("popq	%rbp\n")

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
	File.WriteString(code)
	File.WriteString("retq\n")
}

//PopLocalVariables : Given number of variables pops that amount of variables
func PopLocalVariables(count int) {
	File.WriteString("movq	%rbp, %rsp\n")
	File.WriteString("popq	%rbp\n") //restore rbp

	for i := 0; i < count; i++ {
		//pop remaining local variable
		File.WriteString("popq	%rcx\n")
	}

	//Whenever you update number of variables, remember to reset the rbp address to top of local variables!
	File.WriteString("pushq	%rbp\n")
	File.WriteString("movq	%rsp, %rbp\n")
}
