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

	inStatements := func() {
		if currentNode.Type != "statements" {
			currentNode = childNode(currentNode, "statements")
		}
	}

	for _, token := range tokens {
		switch token.Raw {

		case "{":
			switch currentNode.Type {
			case "subroutineDec":
				currentNode = childNode(currentNode, "subroutineBody")
			case "class":
				insertToken(token, currentNode, true)
				currentNode = childNode(currentNode, "subroutineDec")
				continue
			}
		case "}":
			if currentNode.Type == "subroutineDec" {
				currentNode = currentNode.Parent
			}

			if currentNode.Parent != nil && currentNode.Parent.Type == "whileStatement" {
				currentNode = currentNode.Parent

			}

			insertToken(token, currentNode, true)
			currentNode = currentNode.Parent
			continue
		case "(":
			insertToken(token, currentNode, true)

			if currentNode.Type != "whileStatement" {
				nodeType := "expressionList"
				if currentNode.Type == "subroutineDec" {
					nodeType = "parameterList"
				}
				currentNode = childNode(currentNode, nodeType)
			}
			wrapNextInExpression = true
			continue
		case ")":
			wrapNextInExpression = false
			currentNode = currentNode.Parent

			if currentNode.Type == "expression" {
				currentNode = currentNode.Parent
				if currentNode.Type == "expressionList" || currentNode.Type == "parameterList" {
					currentNode = currentNode.Parent
				}
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
		case ";":
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
			default:
				insertToken(token, currentNode, true)
				currentNode = currentNode.Parent
			}
			continue
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
		case "while":
			inStatements()
			currentNode = childNode(currentNode, "whileStatement")
		case "=":
			insertToken(token, currentNode, true)
			currentNode = childNode(currentNode, "expression")
			continue
		}

		if wrapNextInExpression {
			currentNode = childNode(currentNode, "expression")
			wrapNextInExpression = false
		}

		if currentNode.Type == "expression" {
			currentNode = childNode(currentNode, "term")
		} else if currentNode.Type == "term" && (token.Raw == "<" || token.Raw == "+" || token.Raw == "/") {
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
