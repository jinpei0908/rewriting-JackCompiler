package jacktokenizer

import (
	"fmt"
	_ "io/ioutil"
	"log"
	"os"
	_ "strings"
	"testing"
)

var mapOfTokenType = map[TokenTypes]string{
	Keyword:     "keyword",
	Symbol:      "symbol",
	Identifier:  "identifier",
	IntConst:    "integerConst",
	StringConst: "stringConst",
	None:        "none",
}

/*
var mapOfKeyword = map[KeywordTypes]string{
	_class:       "class",
	_constructor: "constructor",
	_function:    "function",
	_method:      "method",
	_field:       "field",
	_static:      "static",
	_var:         "var",
	_int:         "int",
	_char:        "char",
	_boolean:     "boolean",
	_void:        "void",
	_true:        "true",
	_false:       "false",
	_null:        "null",
	_this:        "this",
	_let:         "let",
	_do:          "do",
	_if:          "if",
	_else:        "else",
	_while:       "while",
	_return:      "return"}

func TestTokenizer(t *testing.T) {
	file, err := os.Open("../../ArrayTest/Main.jack")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	tyFile, err := os.Open("../TokenType.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer tyFile.Close()
	b, _ := ioutil.ReadAll(tyFile)
	tyText := string(b)
	tyInfo := strings.Split(tyText, "\n")

	tk := NewTokenizer(file)
	log.Println(tk.GetTokens())

	for _, tyi := range tyInfo {
		if !tk.HasMoreTokens() {
			break
		}
		tk.Advance()
		actual := mapOfTokenType[tk.TokenType()]
		expect := tyi

		switch tk.TokenType() {
		case Keyword:
			fmt.Println(mapOfKeyword[tk.Keyword()])
		case Symbol:
			fmt.Println(tk.Symbol())
		case StringConst:
			fmt.Println(tk.StringVal())
		case IntConst:
			fmt.Println(tk.IntVal())
		case Identifier:
			fmt.Println(tk.Identifier())
		case None:
			fmt.Println("This is not token. Maybe space.")
		default:
			t.Errorf(`\nUnknown token type: %v\n`, tk.GetCurrentToken())
		}
		if actual != expect {
			t.Errorf("\ntoken: %v\nactual: %v\nexpect: %v\n)",
				tk.GetCurrentToken(), actual, expect)
		}
	}
}
*/

func TestCommaOkIdiom(t *testing.T) {
	v, ok := mapOfTokenType[Keyword]
	fmt.Println(v, ok)
}

func TestGetTokens(t *testing.T) {
	file, err := os.Open("../../Square/SquareGame.jack")
	if err != nil {
		log.Fatalln(err)
	}
	tk := NewTokenizer(file)
	fmt.Println(tk.GetTokens())
}
