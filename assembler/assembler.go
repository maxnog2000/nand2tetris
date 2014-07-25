package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	file, _ := os.Open(os.Args[1])

    addressMap := make(map[string]int64)
	populatePredefinedSymbols(addressMap)

    scanner := bufio.NewScanner(file)
	scanner.Split(consumeLine)
    var currLine int64 = 0

	for ;scanner.Scan(); currLine = currLine + 1 {
		switch currText := scanner.Text(); commandType(currText) {
		case "L_COMMAND":
			addressMap[currText[1:len(currText)-1]] = currLine
			currLine = currLine - 1
		}
	}

	var currAddress int64 = 16
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	scanner.Split(consumeLine)
	for scanner.Scan() {
		switch currText := scanner.Text(); commandType(currText) {
		case "A_COMMAND":
			_, parseError := strconv.ParseInt(currText[1:len(currText)], 10, 64)
			if _, exists := addressMap[currText[1:len(currText)]]; parseError != nil && !exists {
				addressMap[currText[1:len(currText)]] = currAddress
				currAddress = currAddress + 1
			}
		}
	}

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	scanner.Split(consumeLine)

	for scanner.Scan() {
		switch currText := scanner.Text(); commandType(currText) {
		case "A_COMMAND":
			fmt.Println("0" + parseACommand(currText[1:len(currText)], addressMap))
		case "C_COMMAND":
			fmt.Println("111" + comp(currText) + dest(currText) + jump(currText))
		}
	}
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

	//Remove all spaces
	token = bytes.Map(func(r rune) (newR rune) {
		if unicode.IsSpace(r) {
			newR = -1
		} else {
			newR = r
		}
		return
	}, token)
	return
}
