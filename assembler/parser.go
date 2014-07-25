package main

import (
	"strconv"
	"strings"
)

func commandType(ASMLine string) (commandType string) {
	if strings.Contains(ASMLine, "@") {
		return "A_COMMAND"
	} else if strings.Contains(ASMLine, "(") {
		return "L_COMMAND"
	} else {
		return "C_COMMAND"
	}
}

func parseACommand(ASMLine string, addressMap map[string]int64) (binaryString string) {
	intString, parseError := strconv.ParseInt(ASMLine, 10, 64)
	if parseError != nil {
		intString = addressMap[ASMLine]
	}
	binaryString = strconv.FormatInt(intString, 2)
	for binaryStringLen := len(binaryString); binaryStringLen < 15; binaryStringLen = binaryStringLen + 1 {
		binaryString = "0" + binaryString
	}

	return
}
