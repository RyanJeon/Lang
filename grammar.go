package main

import "log"

//list of grammars to classify statements

//VariableDeclare : type for variable declaration statement
type VariableDeclare struct {
	Type      string
	Name      string
	Statement []Token
}

//Classifier tree. Leaf of tree contains the class of statement
//Traverse through the tree to find the class of a statement
type classTree struct {
	Type     string
	Children map[string]classTree
}

var classTreeRoot classTree

//GrammarInit : Initializes grammar tree (Prefix of each statements) **Every type of statement has its own unique prefix to help grammar classification**
func GrammarInit() {
	classTreeRoot = classTree{
		Type:     "Root",
		Children: make(map[string]classTree),
	}

	variableDeclaration := []string{"Declaration", "Variable", "Assignment"}
	addToClassTree(classTreeRoot, variableDeclaration, "VariableDeclaration")

	functionDeclaration := []string{"Declaration", "Function", "=>"}
	addToClassTree(classTreeRoot, functionDeclaration, "FunctionDeclaration")

	functionCall := []string{"Function", "("}
	addToClassTree(classTreeRoot, functionCall, "FunctionCall")

	ifStatement := []string{"If"}
	addToClassTree(classTreeRoot, ifStatement, "IfStatement")

	addToClassTree(classTreeRoot, []string{"Return"}, "FunctionReturn")

	addToClassTree(classTreeRoot, []string{"}"}, "EndOf")

}

func checkValue(t Token) bool {
	return string(t.Value) == "(" || string(t.Value) == ")" || string(t.Value) == "=>" || string(t.Value) == "}"
}

//takes string of term type and add it to classifier tree
func addToClassTree(classTreeRoot classTree, statement []string, class string) {

	start := classTreeRoot
	for _, t := range statement {
		_, exist := start.Children[t]

		if !exist {
			start.Children[t] = classTree{
				Type:     t,
				Children: make(map[string]classTree),
			}
		}
		start = start.Children[t]
	}

	start.Children["StatementType"] = classTree{
		Type:     class,
		Children: make(map[string]classTree),
	}
}

//ClassifyStatement : takes infix token array and classifies the statement
func ClassifyStatement(tokens []Token) string {
	cur := classTreeRoot
	for _, t := range tokens {
		key := t.Type
		//if we should check the value of keyword instead
		if checkValue(t) {
			key = string(t.Value)
		}

		_, exist := cur.Children[key]

		if !exist {
			log.Fatalf("Unexpected %s : (grammar.go Classify)", string(t.Value))
		}

		cur = cur.Children[key]
		statementType, exist := cur.Children["StatementType"]
		if exist {
			return statementType.Type
		}
	}

	if len(cur.Children) > 1 {
		log.Fatalf("Unrecognized statement len(cur.Children) > 1")
	}

	statementType, exist := cur.Children["StatementType"]

	if !exist {
		log.Fatalf("Unrecognized statement : StatementType does not exist")
	}

	return statementType.Type
}
