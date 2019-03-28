package vmwriter

import (
	"fmt"
	"log"
	"os"
)

type VmWriter struct {
	file *os.File
}

func NewVmWriter(outputfile *os.File) *VmWriter {
	return &VmWriter{file: outputfile}
}

func (vm *VmWriter) WritePush(segment string, index int) {
	vm.file.WriteString(
		fmt.Sprintf("push %v %v\n", segment, index))
}

func (vm *VmWriter) WritePop(segment string, index int) {
	vm.file.WriteString(
		fmt.Sprintf("pop %v %v\n", segment, index))
}

func (vm *VmWriter) WriteArithmetic(command string, inTerm bool) {
	switch command {
	case "+":
		vm.file.WriteString("add\n")
	case "-":
		if inTerm {
			vm.file.WriteString("neg\n")
		} else {
			vm.file.WriteString("sub\n")
		}
	case "*":
		vm.file.WriteString("call Math.multiply 2\n")
	case "/":
		vm.file.WriteString("call Math.divide 2\n")
	case "~":
		vm.file.WriteString("not\n")
	case "=":
		vm.file.WriteString("eq\n")
	case "<":
		vm.file.WriteString("lt\n")
	case ">":
		vm.file.WriteString("gt\n")
	case "&":
		vm.file.WriteString("and\n")
	case "|":
		vm.file.WriteString("or\n")
	default:
		log.Fatalln("There is no arithmetic command.")
	}
}

func (vm *VmWriter) WriteLabel(label string) {
	vm.file.WriteString(
		fmt.Sprintf("label %v\n", label))
}

func (vm *VmWriter) WriteGoto(label string) {
	vm.file.WriteString(
		fmt.Sprintf("goto %v\n", label))
}

func (vm *VmWriter) WriteIf(label string) {
	vm.file.WriteString(
		fmt.Sprintf("if-goto %v\n", label))
}

func (vm *VmWriter) WriteCall(name string, nArgs int) {
	vm.file.WriteString(
		fmt.Sprintf("call %v %v\n", name, nArgs))
}

func (vm *VmWriter) WriteFunction(subroutineKind string, className string, subroutineName string, nLocals int, numberOfStatic int) {
	vm.file.WriteString(
		fmt.Sprintf("function %v.%v %v\n", className, subroutineName, nLocals))
	switch subroutineKind {
	case "method":
		vm.file.WriteString("push argument 0\n")
		vm.file.WriteString("pop pointer 0\n")
	case "constructor":
		vm.file.WriteString(fmt.Sprintf("push constant %v\n", numberOfStatic))
		vm.file.WriteString("call Memory.alloc 1\n")
		vm.file.WriteString("pop pointer 0\n")
	}
}

func (vm *VmWriter) WriteReturn() {
	vm.file.WriteString("return\n")
}

func (vm *VmWriter) Close() {
	vm.file.Close()
}
