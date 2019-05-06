package main

import (
	"./compilationengine"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	jackFileNames := []string{}

	arg := getArg(os.Args)

	fInfo, err := os.Stat(arg)
	if err != nil {
		fmt.Println("File information about argument cannot be got")
	}

	if fInfo.IsDir() {
		jackFileNames, err = filepath.Glob(fmt.Sprintf("%v*.jack", arg))
	} else {
		if filepath.Ext(arg) != ".jack" {
			log.Fatalln("Argument is not jack file")
		}
		jackFileNames = append(jackFileNames, arg)
	}

	for _, file := range jackFileNames {
		fmt.Println(file)
		NameOfXML := fmt.Sprintf("%v_.xml", file[:len(file)-5])
		outputXmlFile, err := os.Create(NameOfXML)
		if err != nil {
			log.Fatalln(err)
		}
		defer outputXmlFile.Close()

		inputFile, err := os.Open(file)
		if err != nil {
			log.Fatalln(err)
		}
		defer inputFile.Close()

        outputFile, err := os.Create(fmt.Sprintf("%v_.vm", file[:len(file)-5]))
        if err != nil {
            log.Fatalln(err)
        }
        defer outputFile.Close()

		ce := compilationengine.NewCompilationEngine(inputFile, outputFile, outputXmlFile)
		ce.CompileClass()
	}
}

func getArg(names []string) string {
	if len(names) <= 1 {
		log.Fatalln("Arguments get error: No arg is given")
	} else if len(names) == 2 {
		return names[1]
	} else {
		log.Fatalln("Arguments get error: Too many arguments are given")
	}
	return ""
}
