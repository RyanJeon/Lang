package main

import "log"

//Stack implementation
type Stack []Token

//Queue implementation
type Queue []Token

//Logics for Data Structure starts here !!

//Stack Logics
func (s Stack) Push(v Token) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, Token) {
	l := len(s)
	if l == 0 {
		log.Fatal("Stack is Empty")
	}
	return s[:l-1], s[l-1]
}

func (s Stack) isEmpty() bool {
	return len(s) == 0
}

//Queue Logics
func (q Queue) Add(v Token) Queue {
	return append(q, v)
}

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
