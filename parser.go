package main

import (
	"log"
)

//Tree node
type Tree struct {
	Type  string
	Value []byte
	Left  *Tree
	Right *Tree
}

// * and / have higher precendent than + or -
func precedent(operator string) int {
	if operator == "*" || operator == "/" {
		return 3
	} else if operator == "+" || operator == "-" {
		return 2
	} else {
		log.Fatalf("Error: precedent function takes operator string input: %s", operator)
	}
	return 0
}

//Arithmatic Convert infix to post fix Shunting Yard algo
func TokensPostfix(tokens []Token) []Token {
	operators := make(Stack, 0)
	output := make(Queue, 0)
	postfix := []Token{}

	for _, t := range tokens {
		if t.Type == "Int" || t.Type == "Variable" || t.Type == "FunctionCall" {
			output = output.Add(t)
		} else if t.Type == "Function" ||
			t.Type == "Assignment" ||
			t.Type == "Declaration" ||
			t.Type == "CurlyLeft" {
			operators = operators.Push(t)
		} else if t.Type == "Operator" && t.Value[0] != 40 && t.Value[0] != 41 {
			//While the operator stack is not empty
			for !operators.isEmpty() &&
				(operators.Top().Type == "Function" ||
					(operators.Top().Type == "Operator" &&
						operators.Top().Value[0] != 40 &&
						(precedent(string(operators.Top().Value)) >= precedent(string(t.Value))))) {
				var operator Token
				operators, operator = operators.Pop()
				output = output.Add(operator)
			}
			operators = operators.Push(t)
		}

		//If left paren
		if t.Value[0] == 40 {
			operators = operators.Push(t)
		}
		//If right paren
		if t.Value[0] == 41 {
			//Add to all terms in the parent
			for operators.Top().Value[0] != 40 {
				var operator Token
				operators, operator = operators.Pop()
				output = output.Add(operator)
			}
			//This is empty paren. Discard
			if !operators.isEmpty() && operators.Top().Value[0] == 40 {
				operators, _ = operators.Pop()
			}

		}

	}

	//While the operator stack is not empty
	for !operators.isEmpty() {
		var operator Token
		operators, operator = operators.Pop()
		output = output.Add(operator)
	}

	for !output.isEmpty() {
		var term Token
		output, term = output.Poll()
		postfix = append(postfix, term)
	}

	log.Println(postfix)
	return postfix
}

//Postfix to Abstract Syntax Tree
func tree(post []Token) Tree {
	stack := make(TreeStack, 0)

	for _, t := range post {
		if stack.isEmpty() || t.Type == "Int" || t.Type == "Variable" || t.Type == "FunctionCall" {
			stack = stack.Push(Tree{
				t.Type,
				t.Value,
				nil,
				nil,
			})
			//Type function is temporary for print
		} else if t.Type == "Function" ||
			t.Type == "Declaration" {
			var t1 Tree
			stack, t1 = stack.Pop()
			root := Tree{
				t.Type,
				t.Value,
				nil,
				&t1,
			}
			stack = stack.Push(root)
		} else {
			var t1 Tree
			var t2 Tree
			var root Tree
			stack, t1 = stack.Pop()

			//Could have been a function call before
			stack, t2 = stack.Pop()
			root = Tree{
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

//In Order Traversal for logging tree content / helper to see the tree content
func inorder(tree *Tree) {
	if tree != nil {
		inorder(tree.Left)
		log.Println(string((*tree).Value))
		inorder(tree.Right)
	} else {
		log.Println("NIL")
	}
}
