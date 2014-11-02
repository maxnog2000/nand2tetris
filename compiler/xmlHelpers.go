package main

import (
	"html"
	"io/ioutil"
	"os"

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

func fileToXML(file, baseDirectory string) {
	fileInfo, _ := os.Stat(baseDirectory + file + ".jack")
	xmlOutput := "<tokens>\n"
	for _, token := range JackTokenizer.Tokenize(baseDirectory + file + ".jack") {
		xmlOutput += "<" + tokenTypeFormatter(token.TokenType) + "> " + html.EscapeString(token.Raw) + " </" + tokenTypeFormatter(token.TokenType) + ">\n"
	}
	xmlOutput += "</tokens>\n"
	ioutil.WriteFile(baseDirectory+file+".xml", []byte(xmlOutput), fileInfo.Mode())
}
