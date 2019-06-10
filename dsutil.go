package main

import "log"

//Stack implementation
type Stack []Token

//TreeStack : Stack for Tree Nodes
type TreeStack []Tree

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
		log.Fatal("Stack is Empty")
	}
	return s[:l-1], s[l-1]
}

func (s TreeStack) isEmpty() bool {
	return len(s) == 0
}
