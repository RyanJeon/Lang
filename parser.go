package main

import "log"

type Tree struct {
	Type  string
	Value []byte
	Left  *Tree
	Right *Tree
}

//Arithmatic Convert infix to post fix
func postfix(tokens []Token) []Token {
	operators := make(Stack, 0)
	operands := make(Queue, 0)
	postfix := []Token{}

	for _, t := range tokens {
		if t.Type == "Int" {
			operands = operands.Add(t)
		}
		if t.Type == "Operator" {
			//While the operator stack is not empty
			for !operators.isEmpty() {
				var operator Token
				operators, operator = operators.Pop()
				operands = operands.Add(operator)
			}
			operators = operators.Push(t)
		}
	}

	//While the operator stack is not empty
	for !operators.isEmpty() {
		var operator Token
		operators, operator = operators.Pop()
		operands = operands.Add(operator)
	}

	for !operands.isEmpty() {
		var term Token
		operands, term = operands.Poll()
		postfix = append(postfix, term)
	}

	return postfix
}

//Postfix to Tree
func tree(post []Token) Tree {
	stack := make(TreeStack, 0)

	for _, t := range post {
		if t.Type != "Operator" {
			stack = stack.Push(Tree{
				t.Type,
				t.Value,
				nil,
				nil,
			})
		} else {
			var t1 Tree
			var t2 Tree
			stack, t1 = stack.Pop()
			stack, t2 = stack.Pop()
			root := Tree{
				t.Type,
				t.Value,
				&t2,
				&t1,
			}
			stack = stack.Push(root)
		}
	}

	//Root of the tree
	stack, t := stack.Pop()
	return t
}

//In Order Traversal
func inorder(tree *Tree) {
	if tree != nil {
		inorder(tree.Left)
		log.Println((*tree).Value)
		inorder(tree.Right)
	}
}
