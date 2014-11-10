package main

import (
	"fmt"
	"html"
	"strings"

	"./CompilationEngine"
	"./JackTokenizer"
)

func tokenTypeFormatter(token string) string {
	switch token {
	case "IDENTIFIER":
		return "identifier"
	case "SYMBOL":
		return "symbol"
	case "STRING_CONST":
		return "stringConstant"
	case "KEYWORD":
		return "keyword"
	case "INT_CONST":
		return "integerConstant"
	default:
		return "FUBAR"
	}
}

func nodeTraverse(node *CompilationEngine.Node, level int) {
	padding := strings.Repeat(" ", level*2)
	if node.Terminal == true {
		fmt.Printf(padding+"<%s> %s </%s>\n", tokenTypeFormatter(node.Type), html.EscapeString(node.Value), tokenTypeFormatter(node.Type))
	} else {
		fmt.Printf(padding+"<%s>\n", node.Type)
		for _, node := range node.Children {
			nodeTraverse(node, level+1)
		}
		fmt.Printf(padding+"</%s>\n", node.Type)
	}
}

func fileToXML(file, baseDirectory string) {
	nodeTraverse(CompilationEngine.CompilationEngine(JackTokenizer.Tokenize(baseDirectory+file+".jack")), 0)

	//for _, node := range node.Children[3].Children {
	//	fmt.Printf("node %+v\n", node)
	//}

	//fileInfo, _ := os.Stat(baseDirectory + file + ".jack")
	//xmlOutput := "<tokens>\n"
	//for _, token := range JackTokenizer.Tokenize(baseDirectory + file + ".jack") {
	//	fmt.Printf("token %+v\n", token)
	//	//xmlOutput += "<" + tokenTypeFormatter(token.TokenType) + "> " + html.EscapeString(token.Raw) + " </" + tokenTypeFormatter(token.TokenType) + ">\n"
	//}
	//xmlOutput += "</tokens>\n"
	//ioutil.WriteFile(baseDirectory+file+".xml", []byte(xmlOutput), fileInfo.Mode())
}
