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
				insertToken(token, currentNode, true)
				continue
			case "class":
				insertToken(token, currentNode, true)
				currentNode = childNode(currentNode, "subroutineDec")
				continue
			}
		case "}":
			if currentNode.Type == "subroutineDec" {
				currentNode = currentNode.Parent

			}
			insertToken(token, currentNode, true)
			currentNode = currentNode.Parent
			continue
		case "(":
			nodeType := "parameterList"
			if currentNode.Type == "doStatement" {
				nodeType = "expressionList"
			}
			insertToken(token, currentNode, true)
			currentNode = childNode(currentNode, nodeType)
			continue
		case ")":
			currentNode = currentNode.Parent
		case ";":
			switch currentNode.Type {
			case "expression":
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
		case "=":
			insertToken(token, currentNode, true)
			currentNode = childNode(currentNode, "expression")
			continue
		}

		if currentNode.Type == "expression" {
			currentNode = childNode(currentNode, "term")
			insertToken(token, currentNode, true)
			currentNode = currentNode.Parent
		} else {
			insertToken(token, currentNode, true)

		}

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
