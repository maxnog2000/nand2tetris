package CompilationEngine

import "../JackTokenizer"

type Node struct {
	Type, Value string
	Terminal    bool
	Parent      *Node
	Children    []*Node
}

func CompilationEngine(tokens []JackTokenizer.Token) *Node {
	topNode := &Node{
		Type: "class",
	}
	currentNode := topNode

	wrapNextInExpression := false
	wrapNextInTerm := false

	unaryOpLevel := 0

	inStatements := func() {
		if currentNode.Type != "statements" {
			currentNode = childNode(currentNode, "statements")
		}
	}

	for index, token := range tokens {

		if wrapNextInTerm {
			currentNode = childNode(currentNode, "term")
			wrapNextInTerm = false
		}

		switch token.Raw {
		case "{":
			switch currentNode.Type {
			case "subroutineDec":
				currentNode = childNode(currentNode, "subroutineBody")
			}
		case "}":
			if currentNode.Type == "statements" && currentNode.Parent.Type == "subroutineBody" {
				currentNode = currentNode.Parent

				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent
				currentNode = currentNode.Parent

			} else if currentNode.Type == "class" {
				insertToken(token, currentNode, true)
			} else {

				if currentNode.Type == "subroutineDec" || (currentNode.Parent != nil && (currentNode.Parent.Type == "whileStatement" || currentNode.Parent.Type == "ifStatement")) {
					currentNode = currentNode.Parent
				}

				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent

				if currentNode.Type == "subroutineDec" {
					currentNode = currentNode.Parent
				}

			}
			continue
		case "(":
			if wrapNextInExpression {
				currentNode = childNode(currentNode, "expression")
				currentNode = childNode(currentNode, "term")
			}

			insertToken(token, currentNode, true)
			wrapNextInExpression = true

			if currentNode.Type != "whileStatement" && currentNode.Type != "ifStatement" && !(currentNode.Type == "term" && currentNode.Parent.Parent.Type != "letStatement") {
				nodeType := "expressionList"
				if currentNode.Type == "subroutineDec" {
					nodeType = "parameterList"
					wrapNextInExpression = false
				}

				currentNode = childNode(currentNode, nodeType)
			}
			continue
		case ")":
			wrapNextInExpression = false
			currentNode = currentNode.Parent

			if unaryOpLevel != 0 && currentNode.Parent.Parent.Type != "term" {
				for ; unaryOpLevel > 0; unaryOpLevel = unaryOpLevel - 1 {
					currentNode = currentNode.Parent
				}
			}

			if currentNode.Type == "expression" {
				currentNode = currentNode.Parent
			}
			if currentNode.Type == "expressionList" || currentNode.Type == "parameterList" {
				currentNode = currentNode.Parent
			}

		case "[":
			insertToken(token, currentNode, true)
			wrapNextInExpression = true
			continue
		case "]":
			wrapNextInExpression = false
			currentNode = currentNode.Parent

			if currentNode.Type == "expression" {
				currentNode = currentNode.Parent
			}
		case "=":
			if currentNode.Type != "term" {
				insertToken(token, currentNode, true)
				currentNode = childNode(currentNode, "expression")
			} else {
				currentNode = currentNode.Parent
				insertToken(token, currentNode, true)
				wrapNextInTerm = true
			}
			continue
		case "&":
			if len(tokens) > index && tokens[index+1].Raw != "&" {
				wrapNextInTerm = true
			}
		case ",":
			if currentNode.Type == "term" {
				currentNode = currentNode.Parent
				currentNode = currentNode.Parent
				insertToken(token, currentNode, true)
				wrapNextInExpression = true
				continue
			}
		case ";":
			unaryOpLevel = 0
			switch currentNode.Type {
			case "term":
				currentNode = currentNode.Parent
				currentNode = currentNode.Parent
				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent
			case "returnStatement":
				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent
				currentNode = currentNode.Parent
				wrapNextInExpression = false
			default:
				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent
			}
			continue
		case "-", "~":
			wrapNextInTerm = true
			unaryOpLevel = unaryOpLevel + 1
		case "method", "function", "constructor":
			currentNode = childNode(currentNode, "subroutineDec")
		case "if":
			inStatements()
			currentNode = childNode(currentNode, "ifStatement")
		case "var":
			currentNode = childNode(currentNode, "varDec")
		case "let":
			inStatements()
			currentNode = childNode(currentNode, "letStatement")
		case "do":
			inStatements()
			currentNode = childNode(currentNode, "doStatement")
		case "return":
			inStatements()
			currentNode = childNode(currentNode, "returnStatement")
			wrapNextInExpression = true
			insertToken(token, currentNode, true)
			continue
		case "while":
			inStatements()
			currentNode = childNode(currentNode, "whileStatement")
		case "field":
			currentNode = childNode(currentNode, "classVarDec")
		}

		if wrapNextInExpression {
			currentNode = childNode(currentNode, "expression")
			wrapNextInExpression = false
		}

		if currentNode.Type == "expression" {
			currentNode = childNode(currentNode, "term")
		} else if currentNode.Type == "term" && (token.Raw == "<" || token.Raw == "+" || token.Raw == "/" || token.Raw == "&" || token.Raw == ">" || token.Raw == "-") {
			currentNode = currentNode.Parent
		}

		insertToken(token, currentNode, true)

	}

	return topNode
}

func insertToken(token JackTokenizer.Token, node *Node, Terminal bool) {
	node.Children = append(node.Children, &Node{
		Type:     token.TokenType,
		Value:    token.Raw,
		Terminal: Terminal,
		Parent:   node,
	})
}

func childNode(parent *Node, Type string) *Node {
	child := &Node{
		Type:   Type,
		Parent: parent,
	}
	parent.Children = append(parent.Children, child)
	return child
}
