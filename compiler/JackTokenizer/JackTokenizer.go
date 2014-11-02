package JackTokenizer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Token struct {
	Raw       string
	TokenType string //This should be a const
}

func Tokenize(fileName string) (tokens []Token) {
	file, _ := os.Open(fileName)
	scanner := bufio.NewScanner(file)

	tokenHolder := ""
	line := ""
	inStringParse := false
	inComment := false
	for scanner.Scan() {
		line, inComment = stripComments(scanner.Text(), inComment)
		if inComment {
			continue
		}
		for _, val := range line {
			switch char := string(val); char {
			case " ":
				if inStringParse {
					tokenHolder += char
				} else {
					tokens, tokenHolder = maybeBuildToken(tokens, tokenHolder)
				}
			case "\"":
				if inStringParse = !inStringParse; !inStringParse {
					tokens = append(tokens, Token{Raw: tokenHolder, TokenType: "STRING_CONST"})
					tokenHolder = ""
				}
			case "{", "}", "(", ")", "[", "]",
				".", ",", ";", "+", "-", "*",
				"/", "&", "|", "<", ">", "=", "~":
				tokens, tokenHolder = maybeBuildToken(tokens, tokenHolder)
				tokens = append(tokens, Token{Raw: char, TokenType: "SYMBOL"})
			default:
				tokenHolder += char
			}
		}
		tokens, tokenHolder = maybeBuildToken(tokens, tokenHolder)
	}
	return
}

func stripComments(line string, inComment bool) (string, bool) {

	if split := strings.Split(line, "/*"); len(split) > 1 {
		line = split[0]
		inComment = true
		if split := strings.Split(split[1], "*/"); len(split) > 1 {
			inComment = false
		}

	}
	if split := strings.Split(line, "*/"); len(split) > 1 {
		line = split[1]
		inComment = false
	}

	if split := strings.Split(line, "//"); len(split) > 1 {
		line = split[0]
	}

	return line, inComment
}

func maybeBuildToken(tokens []Token, tokenHolder string) ([]Token, string) {
	tokenHolder = strings.TrimSpace(tokenHolder)

	if tokenHolder != "" {
		tokenType := ""
		switch tokenHolder {
		case "class", "constructor", "function", "method",
			"field", "static", "var", "int", "char", "boolean",
			"void", "true", "false", "null", "this",
			"let", "do", "if", "else", "while", "return":
			tokenType = "KEYWORD"
		default:
			if _, err := strconv.ParseInt(tokenHolder, 10, 0); err == nil {
				tokenType = "INT_CONST"
			} else {
				tokenType = "IDENTIFIER"
			}

		}
		tokens = append(tokens, Token{Raw: tokenHolder, TokenType: tokenType})
		tokenHolder = ""
	}
	return tokens, tokenHolder
}
