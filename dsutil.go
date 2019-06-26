package main

import "log"

//Stack implementation
type Stack []Token

//TreeStack : Stack for Tree Nodes
type TreeStack []Tree

//CallStack : Stack for function calls
type CallStack []Call

//Queue implementation
type Queue []Token

//Logics for Data Structure starts here !!

//Stack Logics

//Push : pushes token to stack
func (s Stack) Push(v Token) Stack {
	return append(s, v)
}

//Pop : returns and pop the top of token stack
func (s Stack) Pop() (Stack, Token) {
	l := len(s)
	if l == 0 {
		log.Fatal("Stack is Empty")
	}
	return s[:l-1], s[l-1]
}

//Top : returns top of the stack
func (s Stack) Top() Token {
	l := len(s)
	if l == 0 {
		log.Fatal("Stack is Empty")
	}
	return s[l-1]
}

func (s Stack) isEmpty() bool {
	return len(s) == 0
}

//Push : pushes token to stack
func (s CallStack) Push(v Call) CallStack {
	return append(s, v)
}

//Pop : returns and pop the top of token stack
func (s CallStack) Pop() (CallStack, Call) {
	l := len(s)
	if l == 0 {
		log.Fatal("CallStack is Empty")
	}
	return s[:l-1], s[l-1]
}

//Top : returns top of the stack
func (s CallStack) Top() Call {
	l := len(s)
	if l == 0 {
		log.Fatal("CallStack is Empty")
	}
	return s[l-1]
}

func (s CallStack) isEmpty() bool {
	return len(s) == 0
}

//Queue Logics

//Add : Enqueue new token to queue
func (q Queue) Add(v Token) Queue {
	return append(q, v)
}

//Poll : Dequeue token from queue and return
func (q Queue) Poll() (Queue, Token) {
	l := len(q)
	if l == 0 {
		log.Fatal("Queue is Empty")
	}
	return q[1:l], q[0]
}

func (q Queue) isEmpty() bool {
	return len(q) == 0
}

//Push for tree node stack
func (s TreeStack) Push(v Tree) TreeStack {
	return append(s, v)
}

//Pop for tree node stack
func (s TreeStack) Pop() (TreeStack, Tree) {
	l := len(s)
	if l == 0 {
		log.Fatal("TreeStack is Empty")
	}
	return s[:l-1], s[l-1]
}

func (s TreeStack) isEmpty() bool {
	return len(s) == 0
}

func reverseByteArray(byteArray []byte) []byte {
	for i := 0; i < len(byteArray)/2; i++ {
		temp := byteArray[i]
		byteArray[i] = byteArray[len(byteArray)-1-i]
		byteArray[len(byteArray)-1-i] = temp
	}

	return byteArray
}

//TokenPrettyLog : logs token with its values in string
func TokenPrettyLog(tokenArray []Token) {
}
