package main

//Convert int ascii value to int
func BytesToInt(bytes []byte) int {
	//0 is 0x30
	zero := 48
	res := 0
	for _, c := range bytes {
		res = res*10 + (int(c) - zero)
	}
	return res
}

func interpret(tree *Tree) int {
	if (*tree).Type != "Operator" {
		return BytesToInt((*tree).Value)
	} else {
		//Plus
		if (*tree).Value[0] == 43 {
			return interpret(tree.Left) + interpret(tree.Right)
		} else if (*tree).Value[0] == 45 {
			return interpret(tree.Left) - interpret(tree.Right)
		} else {
			return 0
		}
	}
}
