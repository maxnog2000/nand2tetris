package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var randLblCount int = 0
var staticGrowth int = 0
var staticOffset int = 16

func getRandLabel() (tmpLabel string) {
	randLblCount = randLblCount + 1
	return "LABEL_RAND" + strconv.Itoa(randLblCount)
}

func writeInit(file *os.File) {
	writeCommands(file, "@256", "D=A", "@SP", "M=D")
	writeCall(file, "Sys.init", "0")
}

func writeCall(file *os.File, functionName string, numArgs string) {
	copyRegisterToStack := func(register string) {
		writeCommands(file, "@"+register, "D=M")
		setAtoSP(file)
		writeCommands(file, "M=D")
		incrementRegister(file, "SP")
	}
	//push return-address
	randLabel := getRandLabel()
	writeCommands(file, "@"+randLabel, "D=A")
	setAtoSP(file)
	writeCommands(file, "M=D")
	incrementRegister(file, "SP")
	//push LCL
	copyRegisterToStack("LCL")
	//push ARG
	copyRegisterToStack("ARG")
	//push THIS
	copyRegisterToStack("THIS")
	//push THAT
	copyRegisterToStack("THAT")
	//ARG = SP - n - 5
	writeCommands(file, "@SP", "D=M", "@"+numArgs, "D=D-A", "@5", "D=D-A", "@ARG", "M=D")
	//LCL = SP
	writeCommands(file, "@SP", "D=M", "@LCL", "M=D")
	//goto f
	writeGoto(file, functionName)
	//(return-address)
	writeLabel(file, randLabel)

}

func writeFunction(file *os.File, functionName string, numArgs string) {
	numArgsInt, parseError := strconv.ParseInt(numArgs, 10, 64)
	if parseError != nil {
		log.Fatal(parseError)
	}
	writeLabel(file, functionName)
	var i int64 = 0
	for ; i < numArgsInt; i++ {
		writePushPop(file, "C_PUSH", "constant", "0")
	}
}

func writeReturn(file *os.File) {
	resetCallerVariable := func(frameDistance string, segment string) {
		writeCommands(file, "@"+frameDistance, "D=A", "@R5", "A=M", "A=A-D", "D=M", "@"+segment, "M=D")
	}
	//FRAME = LCL
	writeCommands(file, "@LCL", "D=M", "@R5", "M=D")
	//RET = *(FRAME - 5)
	writeCommands(file, "@5", "D=A", "@R5", "A=M", "A=A-D", "D=M", "@R6", "M=D")
	//*ARG = pop()
	decrementRegister(file, "SP")
	setAtoSP(file)
	writeCommands(file, "D=M", "@ARG", "A=M", "M=D")
	//SP = ARG+1
	writeCommands(file, "@ARG", "D=M+1", "@SP", "M=D")
	//THAT = *(FRAME - 1)
	resetCallerVariable("1", "THAT")
	//THIS = *(FRAME - 2)
	resetCallerVariable("2", "THIS")
	//ARG = *(FRAME - 3)
	resetCallerVariable("3", "ARG")
	//LCL = *(FRAME - 4)
	resetCallerVariable("4", "LCL")
	//goto RET
	writeCommands(file, "@R6", "A=M", "0;JMP")
}

func writeLabel(file *os.File, label string) {
	writeCommands(file, "("+label+")")
}

func writeIf(file *os.File, label string) {
	decrementRegister(file, "SP")
	setAtoSP(file)
	writeCommands(file, "D=M", "@"+label, "D;JNE")
}

func writeGoto(file *os.File, label string) {
	writeCommands(file, "@"+label, "0;JMP")
}

func writeArithmetic(file *os.File, command string) {
	loadStackTop2 := func() {
		decrementRegister(file, "SP")
		setAtoSP(file)
		writeCommands(file, "D=M")
		decrementRegister(file, "SP")
		setAtoSP(file)
	}
	label1 := getRandLabel()
	label2 := getRandLabel()
	switch command {
	case "add":
		loadStackTop2()
		writeCommands(file, "M=M+D")
	case "sub":
		loadStackTop2()
		writeCommands(file, "M=M-D")
	case "and":
		loadStackTop2()
		writeCommands(file, "M=M&D")
	case "or":
		loadStackTop2()
		writeCommands(file, "M=M|D")
	case "not":
		decrementRegister(file, "SP")
		setAtoSP(file)
		writeCommands(file, "M=!M")
	case "neg":
		decrementRegister(file, "SP")
		setAtoSP(file)
		writeCommands(file, "M=-M")
	case "eq":
		loadStackTop2()
		writeCommands(file, "D=M-D")
		writeCommands(file, "@"+label1, "D;JNE", "@1", "D=-A")
		setAtoSP(file)
		writeCommands(file, "M=D", "@"+label2, "0;JMP", "("+label1+")", "@0", "D=A")
		setAtoSP(file)
		writeCommands(file, "M=D", "("+label2+")")
	case "lt":
		loadStackTop2()
		writeCommands(file, "D=D-M")
		writeCommands(file, "@"+label1, "D;JLE", "@1", "D=-A")
		setAtoSP(file)
		writeCommands(file, "M=D", "@"+label2, "0;JMP", "("+label1+")", "@0", "D=A")
		setAtoSP(file)
		writeCommands(file, "M=D", "("+label2+")")
	case "gt":
		loadStackTop2()
		writeCommands(file, "D=D-M")
		writeCommands(file, "@"+label1, "D;JGE", "@1", "D=-A")
		setAtoSP(file)
		writeCommands(file, "M=D", "@"+label2, "0;JMP", "("+label1+")", "@0", "D=A")
		setAtoSP(file)
		writeCommands(file, "M=D", "("+label2+")")
	}
	incrementRegister(file, "SP")
}

func writePushPop(file *os.File, command string, segment string, index string) {
	if command == "C_POP" {
		decrementRegister(file, "SP")
	}

	if segment == "constant" {
		writeCommands(file, "@"+index, "D=A", "@SP", "A=M", "M=D")
	} else if segment == "pointer" {
		switch index {
		case "0":
			segment = "THIS"
		case "1":
			segment = "THAT"
		}
		if command == "C_PUSH" {
			writeCommands(file, "@"+segment, "D=M")
			setAtoSP(file)
			writeCommands(file, "M=D")
		} else {
			setAtoSP(file)
			writeCommands(file, "D=M", "@"+segment, "M=D")
		}
	} else if segment == "temp" || segment == "static" {
		if segment == "temp" {
			segment = "5"
		} else {
			segment = strconv.Itoa(staticOffset)
			if command != "C_PUSH" {
				staticGrowth = staticGrowth + 1
			}
		}
		if command == "C_PUSH" {
			writeCommands(file, "@"+index, "D=A", "@"+segment, "A=A+D", "D=M")
			setAtoSP(file)
			writeCommands(file, "M=D")
		} else {
			writeCommands(file, "@"+index, "D=A", "@"+segment, "D=A+D", "@R13", "M=D")
			setAtoSP(file)
			writeCommands(file, "D=M", "@13", "A=M", "M=D")
		}
	} else {
		switch segment {
		case "local":
			segment = "LCL"
		case "argument":
			segment = "ARG"
		case "this", "that":
			segment = strings.ToUpper(segment)
		}
		if command == "C_PUSH" {
			writeCommands(file, "@"+index, "D=A", "@"+segment, "A=M+D", "D=M")
			setAtoSP(file)
			writeCommands(file, "M=D")
		} else {
			writeCommands(file, "@"+index, "D=A", "@"+segment, "D=M+D", "@R13", "M=D")
			setAtoSP(file)
			writeCommands(file, "D=M", "@13", "A=M", "M=D")
		}
	}
	if command == "C_PUSH" {
		incrementRegister(file, "SP")
	}
}

func setAtoSP(file *os.File) {
	writeCommands(file, "@SP", "A=M")
}

func incrementRegister(file *os.File, register string) {
	writeCommands(file, "@"+register, "M=M+1")
}

func decrementRegister(file *os.File, register string) {
	writeCommands(file, "@"+register, "M=M-1")
}

func writeCommands(file *os.File, args ...string) {
	for _, arg := range args {
		file.WriteString(arg + "\n")
	}
}
