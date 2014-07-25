package main

import (
	"bufio"
	"os"
	"strings"
)

func compileVM(srcFile *os.File, dstFile *os.File) {
	srcScanner := bufio.NewScanner(srcFile)
	srcScanner.Split(consumeLine)
	staticOffset = staticOffset + staticGrowth
	staticGrowth = 0

	for srcScanner.Scan() {
		writeCommands(dstFile, "//begin "+srcScanner.Text())
		commandType := commandType(srcScanner.Text())
		switch commandType {
		case "C_PUSH", "C_POP":
			writePushPop(dstFile, commandType, arg1(srcScanner.Text()), arg2(srcScanner.Text()))
		case "C_ARITHMETIC":
			writeArithmetic(dstFile, trimString(srcScanner.Text()))
		case "C_LABEL":
			writeLabel(dstFile, arg1(srcScanner.Text()))
		case "C_IF":
			writeIf(dstFile, arg1(srcScanner.Text()))
		case "C_GOTO":
			writeGoto(dstFile, arg1(srcScanner.Text()))
		case "C_RETURN":
			writeReturn(dstFile)
		case "C_FUNCTION":
			writeFunction(dstFile, arg1(srcScanner.Text()), arg2(srcScanner.Text()))
		case "C_CALL":
			writeCall(dstFile, arg1(srcScanner.Text()), arg2(srcScanner.Text()))
		}
		writeCommands(dstFile, "//end "+srcScanner.Text())

	}
}

func trimString(VMLine string) (cleanString string) {
	return strings.Trim(VMLine, " ")
}

func arg1(VMLine string) (arg1 string) {
	return strings.Split(VMLine, " ")[1]
}

func arg2(VMLine string) (arg2 string) {
	return strings.Split(VMLine, " ")[2]
}

func commandType(VMLine string) (commandType string) {
	switch strings.Split(VMLine, " ")[0] {
	case "push":
		return "C_PUSH"
	case "pop":
		return "C_POP"
	case "add", "lt", "neg", "sub", "eq", "or", "not", "gt", "and":
		return "C_ARITHMETIC"
	case "label":
		return "C_LABEL"
	case "if-goto":
		return "C_IF"
	case "goto":
		return "C_GOTO"
	case "return":
		return "C_RETURN"
	case "function":
		return "C_FUNCTION"
	case "call":
		return "C_CALL"
	}
	return "NONE"
}

func consumeLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if advance == 0 && len(token) == 0 {
		return
	}
	tokenString := string(token)
	//Skip empty lines
	if advance == 2 {
		advance, token, err = consumeLine(data[advance:len(data)], atEOF)
		advance = advance + 2
	}
	//Drop lines with comments OR drop line ending with comments
	if strings.IndexAny(tokenString, "/") == 0 {
		storedAdvance := advance
		advance, token, err = consumeLine(data[storedAdvance:len(data)], atEOF)
		advance = advance + storedAdvance
	} else if commentIndex := strings.IndexAny(tokenString, "/"); commentIndex != -1 {
		token = token[0:commentIndex]
	}
	return
}
