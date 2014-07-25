package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Specifiy directory or file")
		return
	}
	fileInfo, fileError := os.Stat(os.Args[1])
	dirtyFilesToProcess := []string{}
	var baseDirectory string

	if fileError != nil {
		fmt.Println("%s", fileError)
		return
	} else if fileInfo.IsDir() {
		directory, _ := os.Open(os.Args[1])
		baseDirectory = fileInfo.Name() + "/"
		dirtyFilesToProcess, _ = directory.Readdirnames(0)
	} else {
		dirtyFilesToProcess = []string{fileInfo.Name()}
		if slashIndex := strings.LastIndex(os.Args[1], "/"); slashIndex != -1 {
			baseDirectory = os.Args[1][0:slashIndex] + "/"
		} else {
			baseDirectory = "."
		}
	}

	filesToProcess := []string{}
	for _, value := range dirtyFilesToProcess {
		if strings.HasSuffix(value, ".vm") {
			filesToProcess = append(filesToProcess, strings.TrimSuffix(value, ".vm"))
		}
	}
	if len(filesToProcess) != 0 {
		dstFile, _ := os.Create(strings.Join([]string{baseDirectory, strings.TrimSuffix(fileInfo.Name(), ".vm"), ".asm"}, ""))
		writeInit(dstFile)
		for _, file := range filesToProcess {
			srcFile, _ := os.Open(strings.Join([]string{baseDirectory, file, ".vm"}, ""))
		    compileVM(srcFile, dstFile)
		}
	}
}
