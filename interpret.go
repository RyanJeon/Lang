package main

//BytesToInt : Convert int ascii value to int
func BytesToInt(bytes []byte) int {
	//0 is 0x30
	zero := 48
	res := 0
	for _, c := range bytes {
		res = res*10 + (int(c) - zero)
	}
	return res
}

//Arithmetic : given Arithmetic tree, calculates integer value
func Arithmetic(tree *Tree) int {
	if (*tree).Type != "Operator" {
		return BytesToInt((*tree).Value)
	}
	val := (*tree).Value[0]
	//Plus
	if val == 43 {
		return Arithmetic(tree.Left) + Arithmetic(tree.Right)
		//Minus
	} else if val == 45 {
		return Arithmetic(tree.Left) - Arithmetic(tree.Right)
		//Mult
	} else if val == 42 {
		return Arithmetic(tree.Left) * Arithmetic(tree.Right)
		//Div
	} else if val == 47 {
		return Arithmetic(tree.Left) / Arithmetic(tree.Right)
	} else {
		return 0
	}
}
