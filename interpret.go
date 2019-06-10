package main

import "log"

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

func arithmetic(tree *Tree) int {
	log.Println(*tree)
	if (*tree).Type != "Operator" {
		return BytesToInt((*tree).Value)
	}
	val := (*tree).Value[0]
	//Plus
	if val == 43 {
		return arithmetic(tree.Left) + arithmetic(tree.Right)
		//Minus
	} else if val == 45 {
		return arithmetic(tree.Left) - arithmetic(tree.Right)
		//Mult
	} else if val == 42 {
		return arithmetic(tree.Left) * arithmetic(tree.Right)
		//Div
	} else if val == 47 {
		return arithmetic(tree.Left) / arithmetic(tree.Right)
	} else {
		return 0
	}
}
