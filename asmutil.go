package main

import "os"

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

//IntToHex : writes integer to decimal equivalent of ascii in current file takes rax as input, and outputs to rcx
func IntToHex(f *os.File) {
	//flip the digits before conversion
	f.WriteString("inttohex:\n")
	flipDigits(f)
	//Initialize variables, r9 = 10 (digit iterator), cl = 0, (8*n) bit shift
	f.WriteString("movq	$10, %r9\n")
	f.WriteString("mov		$0, %cl\n")
	f.WriteString("xor		%rbx, %rbx\n")
	f.WriteString("divide:\n")

	//Evaluate next digit, add digit * 16^n for equivalent ascii value
	f.WriteString("xor		%rdx, %rdx\n")
	f.WriteString("divq	%r9\n")         //val / 10
	f.WriteString("pushq	%rax\n")       //Save current int value
	f.WriteString("movq	$0x30, %rax\n") //Move zero char (0x30) to rax
	f.WriteString("addq	%rdx, %rax\n")  //add the remainder (next digit) to rax
	f.WriteString("salq	%cl, %rax\n")   // shift 8*n bits to left
	f.WriteString("addq	%rax, %rbx\n")  // add the current digit to rcx
	f.WriteString("add		$8, %cl\n")     // add 8 to r8 (shift 8 more bits next time)
	f.WriteString("popq	%rax\n")        //restore current val
	f.WriteString("test	%rax, %rax\n")  //Loop back if not done
	f.WriteString("jnz		divide\n")

	//When done, add end of line
	f.WriteString("add		$8, %cl\n")
	f.WriteString("movq	$0x0A, %rax\n")
	f.WriteString("salq	%cl, %rax\n")
	f.WriteString("addq	%rax, %rbx\n")

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

func PrintAsm(f *os.File) {
	f.WriteString("printout:\n")
	f.WriteString("movq	$0x2000004, %rax\n")
	f.WriteString("movq	$1, %rdi\n")
	f.WriteString("movq	%rcx, print(%rip)\n")
	f.WriteString("leaq	print(%rip), %rsi\n")
	f.WriteString("movq	$100, %rdx\n")
	f.WriteString("syscall\n")
	f.WriteString("retq\n")
}
