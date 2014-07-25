package main

import (
	"strings"
)

func dest(ASMLine string) (binaryString string) {
	if assignIndex := strings.IndexAny(ASMLine, "="); assignIndex == -1 {
		return "000"
	} else {
		switch ASMLine = ASMLine[0:assignIndex]; ASMLine {
		case "M":
			return "001"
		case "D":
			return "010"
		case "MD":
			return "011"
		case "A":
			return "100"
		case "AM":
			return "101"
		case "AD":
			return "110"
		case "AMD":
			return "111"
		default:
			return "000"
		}
	}
}

func comp(ASMLine string) (binaryString string) {
	if assignIndex := strings.IndexAny(ASMLine, "="); assignIndex != -1 {
		ASMLine = ASMLine[assignIndex+1 : len(ASMLine)]
	}

	if JMPIndex := strings.IndexAny(ASMLine, ";"); JMPIndex != -1 {
		ASMLine = ASMLine[0:JMPIndex]
	}
	switch ASMLine {
	case "0":
		return "0101010"
	case "1":
		return "0111111"
	case "-1":
		return "0111010"
	case "D":
		return "0001100"
	case "A":
		return "0110000"
	case "M":
		return "1110000"
	case "!D":
		return "0001101"
	case "!A":
		return "0110001"
	case "!M":
		return "1110001"
	case "D+1":
		return "0011111"
	case "A+1":
		return "0110111"
	case "M+1":
		return "1110111"
	case "D-1":
		return "0001110"
	case "A-1":
		return "0110010"
	case "M-1":
		return "1110010"
	case "D+A":
		return "0000010"
	case "D+M":
		return "1000010"
	case "D-A":
		return "0010011"
	case "D-M":
		return "1010011"
	case "A-D":
		return "0000111"
	case "M-D":
		return "1000111"
	case "D&A":
		return "0000000"
	case "D&M":
		return "1000000"
	case "D|A":
		return "0010101"
	case "D|M":
		return "1010101"
	default:
		return "PANIC"
	}
}

func jump(ASMLine string) (binaryString string) {
	switch {
	case strings.Contains(ASMLine, "JGT"):
		return "001"
	case strings.Contains(ASMLine, "JEQ"):
		return "010"
	case strings.Contains(ASMLine, "JGE"):
		return "011"
	case strings.Contains(ASMLine, "JLT"):
		return "100"
	case strings.Contains(ASMLine, "JNE"):
		return "101"
	case strings.Contains(ASMLine, "JLE"):
		return "110"
	case strings.Contains(ASMLine, "JMP"):
		return "111"
	default:
		return "000"
	}
}

func populatePredefinedSymbols(addressMap map[string]int64) {
	addressMap["SP"] = 0
	addressMap["LCL"] = 1
	addressMap["ARG"] = 2
	addressMap["THIS"] = 3
	addressMap["THAT"] = 4
	addressMap["R0"] = 0
	addressMap["R1"] = 1
	addressMap["R2"] = 2
	addressMap["R3"] = 3
	addressMap["R4"] = 4
	addressMap["R5"] = 5
	addressMap["R6"] = 6
	addressMap["R7"] = 7
	addressMap["R8"] = 8
	addressMap["R9"] = 9
	addressMap["R10"] = 10
	addressMap["R11"] = 11
	addressMap["R12"] = 12
	addressMap["R13"] = 13
	addressMap["R14"] = 14
	addressMap["R15"] = 15
	addressMap["SCREEN"] = 16384
	addressMap["KBD"] = 24576
}
