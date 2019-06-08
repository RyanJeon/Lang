package main

type Tree struct {
	Type  string
	Value []byte
	Next  *Tree
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
