package compilationengine

import (
	"os"
	"testing"
)

func TestComileTerm(t *testing.T) {
	inputFile, _ := os.Open("Script.jack")
	defer inputFile.Close()
	debugFile, _ := os.Create("outputOfTestCompileTerms.xml")
	defer debugFile.Close()
	outputFile, _ := os.Create("Script_.vm")
	defer outputFile.Close()

	cmplEngn := NewCompilationEngine(inputFile, outputFile, debugFile)
	cmplEngn.CompileClass()
}
