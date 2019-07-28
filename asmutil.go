package main

import (
	"fmt"
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
func ASMInit() {
	File.Write([]byte(".data\n"))
	File.Write([]byte("print: .asciz \"\" \n"))
	File.Write([]byte(".globl execute\n")) //make main visible to linker
	File.Write([]byte(".text\n"))
	intToHex()
	printAsm()
}

//IntToHex : writes integer to decimal equivalent of ascii in current file takes rax as input, and outputs to rcx
func intToHex() {
	//flip the digits before conversion
	File.WriteString("inttostring:\n")
	//Initialize variables, r9 = 10 (digit iterator), cl = 0, (8*n) bit shift
	File.WriteString("movq	$10, %r9\n")
	File.WriteString("xor	%rbx, %rbx\n")

	//When done, add end of line
	File.WriteString("movq	$0xa, %rbx\n")
	File.WriteString("salq	$8, %rbx\n")

	//start of adding digits
	File.WriteString("divide:\n")

	//Evaluate next digit, add digit * 16^n for equivalent ascii value
	File.WriteString("xor		%rdx, %rdx\n")
	File.WriteString("divq	%r9\n")         //val / 10
	File.WriteString("pushq	%rax\n")       //Save current int value
	File.WriteString("movq	$0x30, %rax\n") //Move zero char (0x30) to rax
	File.WriteString("addq	%rdx, %rax\n")  //add the remainder (next digit) to rax
	File.WriteString("salq	$8, %rbx\n")    // shift 8*n bits to left
	File.WriteString("addq	%rax, %rbx\n")  // add the current digit to rcx
	File.WriteString("popq	%rax\n")        //restore current val
	File.WriteString("test	%rax, %rax\n")  //Loop back if not done
	File.WriteString("jnz		divide\n")

	//Clear rcx, and move final ascii value to rcx
	File.WriteString("xor		%rcx, %rcx\n")
	File.WriteString("movq	%rbx, %rcx\n")

	File.WriteString("retq\n")

}

//Takes rax and flips its digits
func flipDigits() {
	File.WriteString("xor		%rcx, %rcx\n") //rcx will store result
	File.WriteString("movq	$10, %r9\n")
	File.WriteString("flipdigit:\n")
	File.WriteString("xor		%rdx, %rdx\n")
	File.WriteString("divq	%r9\n")
	File.WriteString("pushq	%rax\n")      //save remaining value
	File.WriteString("pushq	%rdx\n")      //save remainder
	File.WriteString("movq	%rcx, %rax\n") //Move return value to rax
	File.WriteString("mulq	%r9\n")        //return * 10
	File.WriteString("popq	%rdx\n")       //restore remainder
	File.WriteString("addq	%rdx, %rax\n") //add remainder to the result
	File.WriteString("movq	%rax, %rcx\n") //save the result back to rcx
	File.WriteString("popq	%rax\n")       //restore remaining value
	File.WriteString("test	%rax, %rax\n") //Loop back if not done
	File.WriteString("jnz		flipdigit\n")
	File.WriteString("movq	%rcx, %rax\n") //move the result to rax
}

func printAsm() {
	File.WriteString("printout:\n")
	File.WriteString("movq	$0x2000004, %rax\n")
	File.WriteString("movq	$1, %rdi\n")
	File.WriteString("movq	%rcx, print(%rip)\n")
	File.WriteString("leaq	print(%rip), %rsi\n")
	File.WriteString("movq	$100, %rdx\n")
	File.WriteString("syscall\n")
	File.WriteString("retq\n")
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

//Small Utility functions

//ASMWrite takes string and writes it to designated File
func ASMWrite(code string) {
	File.WriteString(code)
}
