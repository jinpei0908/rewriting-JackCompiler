package compilationengine

import (
	. "../jacktokenizer"
	"../symboltable"
	"../vmwriter"
	"fmt"
	"log"
	"os"
)

type compilationEngine struct {
	tk              *Tokenizer
	vm              *vmwriter.VmWriter
	st              *symboltable.SymbolTable
	in              *os.File
	out             *os.File
	outForDebug     *os.File
	thisClassName   string
	numOfExpression int
}

var segments = map[string]string{
	"Static":   "static",
	"Var":      "local",
	"Field":    "this",
	"Argument": "argument",
}

func NewCompilationEngine(inputFile, outputFile, debugFile *os.File) *compilationEngine {
	tokenizer := NewTokenizer(inputFile)
	symbolTable := symboltable.NewSymbolTable()
	vm := vmwriter.NewVmWriter(outputFile)
	return &compilationEngine{
		tk:              tokenizer,
		vm:              vm,
		st:              symbolTable,
		in:              inputFile,
		out:             outputFile,
		outForDebug:     debugFile,
		thisClassName:   "",
		numOfExpression: 0,
	}
}

func (ce *compilationEngine) CompileClass() {
	ce.writeTag("<class>")
	defer ce.writeTag("</class>")

	ce.writeKeyword()    // "class"
	ce.writeIdentifier() // className
	ce.writeIdentifiersInfo("class", true)
	ce.thisClassName = ce.tk.GetCurrentToken()
	ce.writeSymbol() // "{"

	for {
		if next := ce.CheckNextToken(); next != "static" && next != "field" {
			break
		}
		ce.CompileClassVarDec()
	}
	for {
		if next := ce.CheckNextToken(); next != "constructor" && next != "method" && next != "function" {
			break
		}
		ce.CompileSubroutine()
	}
	ce.writeSymbol() // "}"
	ce.showTableOfClass()
}

func (ce *compilationEngine) CompileClassVarDec() {
	ce.writeTag("<classVarDec>")
	defer ce.writeTag("</classVarDec>")

	ce.st.CurrentKind = ce.CheckNextToken()
	ce.writeKeyword() // ("static" | "field")
	ce.st.CurrentType = ce.CheckNextToken()
	ce.writeType() // type

	ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
	ce.writeIdentifier()              // varName
	ce.writeIdentifiersInfo("", true) // its info
	for {
		if ce.CheckNextToken() == ";" {
			break
		}
		ce.writeSymbol() // ","
		ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
		ce.writeIdentifier()              // varName
		ce.writeIdentifiersInfo("", true) // its info
	}
	ce.writeSymbol() // ";"
}

func (ce *compilationEngine) CompileSubroutine() {
	ce.writeTag("<Dec>")
	defer ce.writeTag("</Dec>")

	ce.writeKeyword()                         // ("constructor" | "function" | "method")
	subroutineKind := ce.tk.GetCurrentToken() // subroutineKind = ("constructor" | "function" | "method")
	ce.st.StartSubroutine(subroutineKind)
	ce.writeType()                              // ("void" | type)
	ce.writeIdentifier()                        // Name
	ce.writeIdentifiersInfo("subroutine", true) // its info
	functionName := ce.tk.GetCurrentToken()     // ce.functionName = subroutineName

	ce.writeSymbol()          // "("
	ce.CompileParameterList() //
	ce.writeSymbol()          // ")"

	ce.writeTag("<Body>")
	defer ce.writeTag("</Body>")

	ce.writeSymbol() // "{"

	for {
		if ce.CheckNextToken() != "var" {
			break
		}
		ce.CompileVarDec()
	}

	numberOfVars, err := ce.st.VarCount("Var")
	if err != nil {
		log.Fatalln(err, "Fail to get number of variables in compileSubroutine")
	}
	numberOfField, err := ce.st.VarCount("Field")
	if err != nil {
		log.Fatalln(err, "Fail to get number of field variables in compileSubroutine")
	}

	ce.vm.WriteFunction(
		subroutineKind,
		ce.thisClassName,
		functionName,
		numberOfVars,
		numberOfField)

	ce.CompileStatements()
	ce.writeSymbol() // "}"
	ce.showTableOfSubroutine()
}

func (ce *compilationEngine) CompileParameterList() {
	ce.writeTag("<parameterList>")
	defer ce.writeTag("</parameterList>")

	if ce.CheckNextToken() == ")" {
		return
	}

	ce.st.CurrentKind = "arg"
	ce.st.CurrentType = ce.CheckNextToken()
	ce.writeType() // type
	ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
	ce.writeIdentifier()              // varName
	ce.writeIdentifiersInfo("", true) // its info

	for {
		if ce.tk.CheckNextToken() != "," {
			break
		}
		ce.writeSymbol() // ","
		ce.st.CurrentKind = "arg"
		ce.st.CurrentType = ce.CheckNextToken()
		ce.writeType() // type
		ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
		ce.writeIdentifier()              // varName
		ce.writeIdentifiersInfo("", true) // its info
	}
}

func (ce *compilationEngine) CompileVarDec() {
	ce.writeTag("<varDec>")
	defer ce.writeTag("</varDec>")

	ce.st.CurrentKind = ce.CheckNextToken()
	ce.writeKeyword() // "var"
	ce.st.CurrentType = ce.CheckNextToken()
	ce.writeType() // type

	ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
	ce.writeIdentifier()              // varName
	ce.writeIdentifiersInfo("", true) // its info

	for {
		if ce.CheckNextToken() == ";" {
			break
		}
		ce.writeSymbol() // ","
		ce.st.Define(ce.CheckNextToken(), ce.st.CurrentType, ce.st.CurrentKind)
		ce.writeIdentifier()              // varName
		ce.writeIdentifiersInfo("", true) // its info
	}
	ce.writeSymbol() // ";"
}

func (ce *compilationEngine) CompileStatements() {
	ce.writeTag("<statements>")
	defer ce.writeTag("</statements>")
	for {
		switch ce.CheckNextToken() {
		case "let":
			ce.CompileLet()
		case "if":
			ce.CompileIf()
		case "while":
			ce.CompileWhile()
		case "do":
			ce.CompileDo()
		case "return":
			ce.CompileReturn()
		default:
			return
		}
	}
}

func (ce *compilationEngine) CompileDo() {
	ce.writeTag("<doStatement>")
	defer ce.writeTag("</doStatement>")

	ce.writeKeyword() // "do"
	ce.compileSubroutineCall()

	ce.writeSymbol() // ";"
	ce.vm.WritePop("temp", 0)
}

func (ce *compilationEngine) CompileLet() {
	ce.writeTag("<letStatement>")
	defer ce.writeTag("</letStatement>")

	ce.writeKeyword()    // "let"
	ce.writeIdentifier() // varName
	varName := ce.tk.Identifier()
	varNameKind, _ := ce.st.KindOf(varName)
	varNameIndex, _ := ce.st.IndexOf(varName)

	ce.writeIdentifiersInfo("", false) // its info

	// If varName is array
	if ce.CheckNextToken() == "[" {
		ce.writeSymbol() // "["
		ce.CompileExpression()
		ce.writeSymbol() // "]"
		ce.vm.WritePush(segments[varNameKind], varNameIndex)
		ce.vm.WriteArithmetic("+", false)

		ce.writeSymbol() // "="
		ce.CompileExpression()

		ce.vm.WritePop("temp", 0)
		ce.vm.WritePop("pointer", 1)
		ce.vm.WritePush("temp", 0)
		ce.vm.WritePop("that", 0)

	} else {
		ce.writeSymbol() // "="
		ce.CompileExpression()
		ce.vm.WritePop(segments[varNameKind], varNameIndex)
	}

	ce.writeSymbol() // ";"
}

func (ce *compilationEngine) CompileWhile() {
	ce.writeTag("<whileStatement>")
	defer ce.writeTag("</whileStatement>")

	whileStart := fmt.Sprintf("WHILE_EXP%v", ce.st.WhileCount)
	whileEnd := fmt.Sprintf("WHILE_END%v", ce.st.WhileCount)
	ce.st.WhileCount++

	ce.vm.WriteLabel(whileStart)

	ce.writeKeyword() // "while"
	ce.writeSymbol()  // "("
	ce.CompileExpression()
	ce.writeSymbol() // ")"

	ce.vm.WriteArithmetic("~", false)
	ce.vm.WriteIf(whileEnd)
	ce.writeSymbol() // "{"
	ce.CompileStatements()
	ce.vm.WriteGoto(whileStart)
	ce.vm.WriteLabel(whileEnd)
	ce.writeSymbol() // "}"
}

func (ce *compilationEngine) CompileReturn() {
	ce.writeTag("<returnStatement>")
	defer ce.writeTag("</returnStatement>")

	ce.writeKeyword() // "return"
	if ce.CheckNextToken() != ";" {
		ce.CompileExpression()
	} else {
		ce.vm.WritePush("constant", 0)
	}
	ce.writeSymbol() // ";"
	ce.vm.WriteReturn()
}

func (ce *compilationEngine) CompileIf() {
	ce.writeTag("<ifStatement>")
	defer ce.writeTag("</ifStatement>")

	trueLabel := fmt.Sprintf("IF_TRUE%v", ce.st.IfCount)
	falseLabel := fmt.Sprintf("IF_FALSE%v", ce.st.IfCount)
	endLabel := fmt.Sprintf("IF_END%v", ce.st.IfCount)
	ce.st.IfCount++

	ce.writeKeyword()      // "if"
	ce.writeSymbol()       // "("
	ce.CompileExpression() // expression
	ce.writeSymbol()       // ")"
	ce.writeSymbol()       // "{"

	ce.vm.WriteIf(trueLabel)
	ce.vm.WriteGoto(falseLabel)
	ce.vm.WriteLabel(trueLabel)
	ce.CompileStatements() // statements
	ce.writeSymbol()       // }

	if ce.CheckNextToken() != "else" {
		ce.vm.WriteLabel(falseLabel)
		return
	}
	ce.vm.WriteGoto(endLabel)

	ce.writeKeyword() // "else"
	ce.writeSymbol()  // "{"
	ce.vm.WriteLabel(falseLabel)
	ce.CompileStatements() // statements
	ce.writeSymbol()       // "}"
	ce.vm.WriteLabel(endLabel)
}

func (ce *compilationEngine) CompileExpression() {
	ce.writeTag("<expression>")
	defer ce.writeTag("</expression>")

	ce.CompileTerm()

	for {
		switch ce.CheckNextToken() {
		case "+", "-", "*", "/", "&", "|", "<", ">", "=":
			// Write op, term
			ce.writeSymbol()
			op := ce.tk.GetCurrentToken()
			ce.CompileTerm()
			ce.vm.WriteArithmetic(op, false)

		default:
			return
		}
	}
}

func (ce *compilationEngine) CompileExpressionList() {
	ce.writeTag("<expressionList>")
	defer ce.writeTag("</expressionList>")

	numOfExpr := 0
	defer func() { ce.numOfExpression = numOfExpr }()
	if ce.CheckNextToken() == ")" {
		return
	}

	ce.CompileExpression()
	numOfExpr++
	fmt.Println(numOfExpr)

	for {
		// No more expression exist, return
		if ce.CheckNextToken() != "," {
			break
		}

		// Write ",", expression
		ce.writeSymbol()
		ce.CompileExpression()
		numOfExpr++
		fmt.Println(numOfExpr)
	}
}

func (ce *compilationEngine) CompileTerm() {
	ce.writeTag("<term>")
	defer ce.writeTag("</term>")

	if !ce.tk.HasMoreTokens() {
		return
	}
	ce.tk.Advance()
	switch tt := ce.tk.TokenType(); tt {
	case IntConst:
		// ce.writeTokenWithTag(strconv.Itoa(ce.tk.IntVal()), "integerConstant")
		ce.vm.WritePush("constant", ce.tk.IntVal())

	case StringConst:
		ce.writeTokenWithTag(ce.tk.StringVal(), "stringConstant")
		str := ce.tk.StringVal()
		strLength := len(str)
		ce.vm.WritePush("constant", strLength)
		ce.vm.WriteCall("String.new", 1)

		for _, char := range str {
			ce.vm.WritePush("constant", int(byte(char)))
			ce.vm.WriteCall("String.appendChar", 2)
		}

	case Keyword:
		ce.writeTokenWithTag(ce.tk.Keyword(), "keyword")
		token := ce.tk.Keyword()
		switch token {
		case "this":
			ce.vm.WritePush("pointer", 0)
		case "that":
			ce.vm.WritePush("pointer", 1)
		case "true":
			ce.vm.WritePush("constant", 0)
			ce.vm.WriteArithmetic("~", true)
		case "false", "null":
			ce.vm.WritePush("constant", 0)
		}

	case Symbol:
		switch token := ce.tk.Symbol(); token {
		case "(":
			// Write "(", expression, ")"
			ce.writeTokenWithTag(ce.tk.Symbol(), "symbol")
			ce.CompileExpression()
			ce.writeSymbol()

		case "-", "~":
			// Write "-", "~"
			ce.writeTokenWithTag(ce.tk.Symbol(), "symbol")
			op := ce.tk.Symbol()
			ce.CompileTerm()
			ce.vm.WriteArithmetic(op, true)
		default:
			log.Fatalln("Failed to get a symbol", token)
		}

	case Identifier:
		// Write identifier
		ce.writeTokenWithTag(ce.tk.Identifier(), "identifier")
		varName := ce.tk.GetCurrentToken()

		switch nextToken := ce.CheckNextToken(); nextToken {
		// Write array
		case "[":
			ce.writeIdentifiersInfo("", false)
			varName := ce.tk.Symbol()
			varNameKind, err := ce.st.KindOf(varName)
			if err != nil {
				log.Fatalln(err, "No kind is found", varName)
			}
			varNameIndex, err := ce.st.IndexOf(varName)
			if err != nil {
				log.Fatalln(err, "No index is found", varName)
			}

			ce.writeSymbol() // Write "["
			ce.CompileExpression()
			ce.writeSymbol() // Write "]"
			ce.vm.WritePush(segments[varNameKind], varNameIndex)
			ce.vm.WriteArithmetic("+", true)
			ce.vm.WritePop("pointer", 1)
			ce.vm.WritePush("that", 0)
			/*
				switch varNameKind {
				case "Static":
					ce.vm.WritePush("static", varNameIndex)
				case "Var":
					ce.vm.WritePush("local", varNameIndex)
				case "Field":
					ce.vm.WritePush("this", varNameIndex)
				case "Argument":
					ce.vm.WritePush("argument", varNameIndex)
				default:
					log.Fatalln("This token is not registered in symbol table")
				}
			*/

		// Write a subroutine in a same class
		case "(":
			ce.writeIdentifiersInfo("subroutine", false)

			functionName := varName
			ce.vm.WritePush("pointer", 0)

			ce.writeSymbol() // "("
			ce.CompileExpressionList()
			ce.writeSymbol() // ")"

			ce.vm.WriteCall(
				fmt.Sprintf("%v.%v", ce.thisClassName, functionName),
				ce.numOfExpression+1)

		// Write a subroutine in a differenct class
		case ".":
			// For symbol table display
			_, err := ce.st.KindOf(varName)
			if err != nil {
				ce.writeIdentifiersInfo("class", false)
			} else {
				ce.writeIdentifiersInfo("", false)
			}
			// end

			ce.writeSymbol()     // "."
			ce.writeIdentifier() // subroutine Name
			ce.writeIdentifiersInfo("subroutine", false)
			subroutineName := ce.tk.Identifier()

			if ce.isInstanceName(varName) {
				typeOfCurrentToken, err := ce.st.TypeOf(varName)
				if err != nil {
					log.Fatalln("Failed to get a type", err)
				}
				kindOfCurrentToken, err := ce.st.KindOf(varName)
				if err != nil {
					log.Fatalln("Failed to get a kind", err)
				}
				indexOfCurrentToken, err := ce.st.IndexOf(varName)
				if err != nil {
					log.Fatalln("Failed to get an index", err)
				}

				ce.vm.WritePush(segments[kindOfCurrentToken], indexOfCurrentToken)

				ce.writeSymbol() // "("
				ce.CompileExpressionList()
				ce.writeSymbol() // ")"

				ce.vm.WriteCall(
					fmt.Sprintf("%v.%v", typeOfCurrentToken, subroutineName),
					ce.numOfExpression+1)

			} else if ce.isClassName(varName) {
				ce.writeSymbol() // "("
				ce.CompileExpressionList()
				ce.writeSymbol() // ")"

				ce.vm.WriteCall(
					fmt.Sprintf("%v.%v", varName, subroutineName),
					ce.numOfExpression)

			} else {
				log.Fatalln("Failed to compile do-statement",
					fmt.Sprintf("%v.%v", varName, subroutineName))
			}
		default:
			ce.writeIdentifiersInfo("", false)

			varName := ce.tk.Identifier()
			varNameKind, err := ce.st.KindOf(varName)
			if err != nil {
				log.Fatalln(err, "No kind is found", varName)
			}
			varNameIndex, err := ce.st.IndexOf(varName)
			if err != nil {
				log.Fatalln(err, "No index is found", varName)
			}

			switch varNameKind {
			case "Static":
				ce.vm.WritePush("static", varNameIndex)
			case "Var":
				ce.vm.WritePush("local", varNameIndex)
			case "Field":
				ce.vm.WritePush("this", varNameIndex)
			case "Argument":
				ce.vm.WritePush("argument", varNameIndex)
			default:
				log.Fatalln("This token is not registered in symbol table")
			}
		}
	}
}

func (ce *compilationEngine) CheckNextToken() string {
	return ce.tk.CheckNextToken()
}

func (ce *compilationEngine) writeTag(s string) {
	ce.outForDebug.WriteString(s + "\n")
}

func (ce *compilationEngine) writeTokenWithTag(s, tagName string) {
	prefixTag := fmt.Sprintf("<%s>", tagName)
	suffixTag := fmt.Sprintf("</%s>", tagName)
	ce.outForDebug.WriteString(
		fmt.Sprintf("%s %s %s\n", prefixTag, s, suffixTag))
}

func (ce *compilationEngine) writeSymbol() {
	if !ce.tk.HasMoreTokens() {
		return
	}
	ce.tk.Advance()
	if ce.tk.TokenType() != Symbol {
		log.Fatalln("Failed to get a symbol", ce.tk.GetCurrentToken())
	}
	ce.writeTokenWithTag(ce.tk.Symbol(), "symbol")
}

func (ce *compilationEngine) writeKeyword() {
	if !ce.tk.HasMoreTokens() {
		return
	}
	ce.tk.Advance()
	if ce.tk.TokenType() != Keyword {
		log.Fatalln("Failed to get a keyword", ce.tk.GetCurrentToken())
	}
	ce.writeTokenWithTag(ce.tk.Keyword(), "keyword")
}

func (ce *compilationEngine) writeIdentifier() {
	if !ce.tk.HasMoreTokens() {
		return
	}
	ce.tk.Advance()
	if ce.tk.TokenType() != Identifier {
		log.Fatalln("Failed to get an identifier", ce.tk.GetCurrentToken())
	}
	ce.writeTokenWithTag(ce.tk.Identifier(), "identifier")
}

func (ce *compilationEngine) writeIdentifiersInfo(category string, define bool) {
	switch category {
	case "class":
		if define {
			ce.writeTokenWithTag("defined class", "identifierInfo")
		} else {
			ce.writeTokenWithTag("used class", "identifierInfo")
		}
	case "subroutine":
		if define {
			ce.writeTokenWithTag("defined subroutine", "identifierInfo")
		} else {
			ce.writeTokenWithTag("used subroutine", "identifierInfo")
		}
	case "":
		if define {
			kind, _ := ce.st.KindOf(ce.tk.Identifier())
			index, _ := ce.st.IndexOf(ce.tk.Identifier())
			type_, _ := ce.st.TypeOf(ce.tk.Identifier())
			info := fmt.Sprintf("defined %v %v %v", kind, index, type_)
			ce.writeTokenWithTag(info, "identifierInfo")
		} else {
			kind, _ := ce.st.KindOf(ce.tk.Identifier())
			index, _ := ce.st.IndexOf(ce.tk.Identifier())
			type_, _ := ce.st.TypeOf(ce.tk.Identifier())
			info := fmt.Sprintf("used %v %v %v", kind, index, type_)
			ce.writeTokenWithTag(info, "identifierInfo")
		}
	}
}

func (ce *compilationEngine) writeType() {
	if !ce.tk.HasMoreTokens() {
		return
	}
	ce.tk.Advance()
	if ce.tk.TokenType() == Keyword {
		ce.writeTokenWithTag(ce.tk.Keyword(), "keyword") // When embeded type
	} else if ce.tk.TokenType() == Identifier {
		ce.writeTokenWithTag(ce.tk.Identifier(), "identifier") // When class
	} else {
		log.Fatalln("Failed to write type", ce.tk.GetCurrentToken())
	}
}

// This method is for debugging.
// When this section completed, this method will be deleted
func (ce compilationEngine) showTableOfSubroutine() {
	fmt.Println("Table of identifiers")
	fmt.Printf("%10v | %10v | %10v | %10v\n", "KEY", "KIND", "TYPE_", "INDEX")
	for key := range ce.st.TableOfSubroutineScope {
		kind, err := ce.st.KindOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		type_, err := ce.st.TypeOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		index, err := ce.st.IndexOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%10v | %10v | %10v | %10v\n", key, kind, type_, index)
	}
}

// This method is for debugging.
// When this section completed, this method will be deleted
func (ce compilationEngine) showTableOfClass() {
	fmt.Println("Table of identifiers")
	fmt.Printf("%10v | %10v | %10v | %10v\n", "KEY", "KIND", "TYPE_", "INDEX")
	for key := range ce.st.TableOfClassScope {
		kind, err := ce.st.KindOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		type_, err := ce.st.TypeOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		index, err := ce.st.IndexOf(key)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%10v | %10v | %10v | %10v\n", key, kind, type_, index)
	}
}

func (ce *compilationEngine) compileSubroutineCall() {
	ce.writeIdentifier()
	currentToken := ce.tk.Identifier()
	nextToken := ce.tk.CheckNextToken()

	switch nextToken {
	// Method in this class
	case "(":
		ce.writeIdentifiersInfo("subroutine", false)
		functionName := currentToken
		ce.vm.WritePush("pointer", 0)

		ce.writeSymbol() // "("
		ce.CompileExpressionList()
		ce.writeSymbol() // ")"

		ce.vm.WriteCall(
			fmt.Sprintf("%v.%v", ce.thisClassName, functionName),
			ce.numOfExpression+1)

	// Function of out of this class or Method
	case ".":
		// For symbol table display
		_, err := ce.st.KindOf(currentToken)
		if err != nil {
			ce.writeIdentifiersInfo("class", false)
		} else {
			ce.writeIdentifiersInfo("", false)
		}
		// end

		ce.writeSymbol()     // "."
		ce.writeIdentifier() // subroutine Name
		ce.writeIdentifiersInfo("subroutine", false)
		subroutineName := ce.tk.Identifier()

		if ce.isInstanceName(currentToken) {
			typeOfCurrentToken, err := ce.st.TypeOf(currentToken)
			if err != nil {
				log.Fatalln("Failed to get a type", err)
			}
			kindOfCurrentToken, err := ce.st.KindOf(currentToken)
			if err != nil {
				log.Fatalln("Failed to get a kind", err)
			}
			indexOfCurrentToken, err := ce.st.IndexOf(currentToken)
			if err != nil {
				log.Fatalln("Failed to get an index", err)
			}

			ce.vm.WritePush(segments[kindOfCurrentToken], indexOfCurrentToken)

			ce.writeSymbol() // "("
			ce.CompileExpressionList()
			ce.writeSymbol() // ")"

			ce.vm.WriteCall(
				fmt.Sprintf("%v.%v", typeOfCurrentToken, subroutineName),
				ce.numOfExpression+1)

		} else if ce.isClassName(currentToken) {
			ce.writeSymbol() // "("
			ce.CompileExpressionList()
			ce.writeSymbol() // ")"

			ce.vm.WriteCall(
				fmt.Sprintf("%v.%v", currentToken, subroutineName),
				ce.numOfExpression)

		} else {
			log.Fatalln("Failed to compile do-statement",
				fmt.Sprintf("%v.%v", currentToken, subroutineName))
		}
	}
}

func (ce *compilationEngine) isInstanceName(id string) bool {
	_, err := ce.st.TypeOf(id)
	if err != nil {
		return false
	}
	return true
}

func (ce *compilationEngine) isClassName(id string) bool {
	_, err := ce.st.TypeOf(id)
	if err == nil {
		return false
	}
	return true
}
