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
	flipDigits(f)

	//Initialize variables, r9 = 10 (digit iterator), cl = 0, (8*n) bit shift
	f.Write([]byte("movq	$10, %r9\n"))
	f.Write([]byte("mov		$0, %cl\n"))
	f.Write([]byte("xor		%rbx, %rbx\n"))
	f.Write([]byte("divide:\n"))

	//Evaluate next digit, add digit * 16^n for equivalent ascii value
	f.Write([]byte("xor		%rdx, %rdx\n"))
	f.Write([]byte("divq	%r9\n"))         //val / 10
	f.Write([]byte("pushq	%rax\n"))       //Save current int value
	f.Write([]byte("movq	$0x30, %rax\n")) //Move zero char (0x30) to rax
	f.Write([]byte("addq	%rdx, %rax\n"))  //add the remainder (next digit) to rax
	f.Write([]byte("salq	%cl, %rax\n"))   // shift 8*n bits to left
	f.Write([]byte("addq	%rax, %rbx\n"))  // add the current digit to rcx
	f.Write([]byte("add		$8, %cl\n"))     // add 8 to r8 (shift 8 more bits next time)
	f.Write([]byte("popq	%rax\n"))        //restore current val
	f.Write([]byte("test	%rax, %rax\n"))  //Loop back if not done
	f.Write([]byte("jnz		divide\n"))

	//When done, add end of line
	f.Write([]byte("add		$8, %cl\n"))
	f.Write([]byte("movq	$0x0A, %rax\n"))
	f.Write([]byte("salq	%cl, %rax\n"))
	f.Write([]byte("addq	%rax, %rbx\n"))

	//Clear rcx, and move final ascii value to rcx
	f.Write([]byte("xor		%rcx, %rcx\n"))
	f.Write([]byte("movq	%rbx, %rcx\n"))

}

//Takes rax and flips its digits
func flipDigits(f *os.File) {
	f.Write([]byte("xor		%rcx, %rcx\n")) //rcx will store result
	f.Write([]byte("movq	$10, %r9\n"))
	f.Write([]byte("flipdigit:\n"))
	f.Write([]byte("xor		%rdx, %rdx\n"))
	f.Write([]byte("divq	%r9\n"))
	f.Write([]byte("pushq	%rax\n"))      //save remaining value
	f.Write([]byte("pushq	%rdx\n"))      //save remainder
	f.Write([]byte("movq	%rcx, %rax\n")) //Move return value to rax
	f.Write([]byte("mulq	%r9\n"))        //return * 10
	f.Write([]byte("popq	%rdx\n"))       //restore remainder
	f.Write([]byte("addq	%rdx, %rax\n")) //add remainder to the result
	f.Write([]byte("movq	%rax, %rcx\n")) //save the result back to rcx
	f.Write([]byte("popq	%rax\n"))       //restore remaining value
	f.Write([]byte("test	%rax, %rax\n")) //Loop back if not done
	f.Write([]byte("jnz		flipdigit\n"))
	f.Write([]byte("movq	%rcx, %rax\n")) //move the result to rax
}
