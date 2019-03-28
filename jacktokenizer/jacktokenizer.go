package jacktokenizer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TokenTypes int
type KeywordTypes int

const (
	Keyword TokenTypes = iota
	Symbol
	Identifier
	IntConst
	StringConst
	None
)

type Tokenizer struct {
	index        int
	tokens       []string
	mapOfString  map[string]string
	currentToken string
}

var patternOfInteger = regexp.MustCompile("[0-9]+")
var patternOfIdentifier = regexp.MustCompile("[A-z_].*")
var symbols = []string{
	"{", "}", "(", ")", "[", "]", ".", ",",
	";", "+", "-", "*", "/", "&",
	"|", "<", ">", "=", "~"}
var keywords = []string{
	"class", "constructor", "function",
	"method", "field", "static", "var",
	"int", "char", "boolean", "void",
	"true", "false", "null", "this",
	"let", "do", "if", "else",
	"while", "return"}

func NewTokenizer(file *os.File) *Tokenizer {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	b, newMap := replaceStringConstant(b)
	b = removeOneLineComment(b)
	b = removeNewLineAndTab(b)
	b = removeMultiLineComment(b)
	s := string(b)
	symbols := strings.NewReplacer(
		"{", " { ",
		"}", " } ",
		"(", " ( ",
		")", " ) ",
		"[", " [ ",
		"]", " ] ",
		".", " . ",
		",", " , ",
		";", " ; ",
		"+", " + ",
		"-", " - ",
		"*", " * ",
		"/", " / ",
		"&", " & ",
		"|", " | ",
		"<", " < ",
		">", " > ",
		"=", " = ",
		"~", " ~ ",
	)

	s = symbols.Replace(s)
	s = replaceMultiSpaceToMonoSpace(s)

	tokens := strings.Split(s, " ")

	return &Tokenizer{
		index:        0,
		tokens:       tokens,
		mapOfString:  newMap,
		currentToken: ""}
}

func (tk *Tokenizer) HasMoreTokens() bool {
	return tk.index < len(tk.tokens)
}

func (tk *Tokenizer) Advance() {
	if tk.tokens[tk.index] == "" {
		tk.index++
	}
	tk.currentToken = tk.tokens[tk.index]
	tk.index++
}

func (tk *Tokenizer) TokenType() TokenTypes {
	if isKeyword(tk.currentToken) {
		return Keyword
	}
	if isSymbol(tk.currentToken) {
		return Symbol
	}
	if strings.HasPrefix(tk.currentToken, "StringConstant_") {
		return StringConst
	}
	if patternOfInteger.MatchString(tk.currentToken) {
		return IntConst
	}
	if patternOfIdentifier.MatchString(tk.currentToken) {
		return Identifier
	}
	return None
}

func (tk *Tokenizer) Keyword() string {
	for _, k := range keywords {
		if tk.currentToken == k {
			return tk.currentToken
		}
	}
	log.Fatalln("This token is not keyword", tk.currentToken)
	return ""
}

func (tk *Tokenizer) Symbol() string {
	switch tk.currentToken {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return tk.currentToken
	}
}

func (tk *Tokenizer) Identifier() string {
	return tk.currentToken
}

func (tk *Tokenizer) IntVal() int {
	i, err := strconv.Atoi(tk.currentToken)
	if err != nil {
		log.Fatalln(err, tk.currentToken)
	}
	return i
}

func (tk *Tokenizer) StringVal() string {
	stringConstant, ok := tk.mapOfString[tk.currentToken]
	if ok {
		return stringConstant
	}
	log.Fatalln("Token is not registered hashmap: ", tk.currentToken)
	return ""
}

func (tk *Tokenizer) CheckNextToken() string {
	if tk.index >= len(tk.tokens) {
		fmt.Println("There is no next token")
		return ""
	}
	i := tk.index
	return tk.tokens[i]
}

func removeMultiLineComment(b []byte) []byte {
	comment := regexp.MustCompile(`/\*.*?\*/`)
	return comment.ReplaceAll(b, []byte(` `))
}

func removeOneLineComment(b []byte) []byte {
	comment := regexp.MustCompile(`//.*\n`)
	return comment.ReplaceAll(b, []byte(` `))
}

func removeNewLineAndTab(b []byte) []byte {
	newLineAndTab := regexp.MustCompile(`\r\n|\r|\n|\t`)
	return newLineAndTab.ReplaceAll(b, []byte(` `))
}

func replaceStringConstant(b []byte) ([]byte, map[string]string) {
	stringConstant := regexp.MustCompile(`"(.*)"`)
	found := stringConstant.FindAllStringSubmatch(string(b), -1)
	mapOfStringConstant := make(map[string]string)

	for i, f := range found {
		s := fmt.Sprintf("StringConstant_%d", i)
		mapOfStringConstant[s] = f[1]
	}

	inputText := string(b)
	for k, v := range mapOfStringConstant {
		s := strings.NewReplacer(v, k)
		inputText = s.Replace(inputText)
	}
	s := strings.NewReplacer(`"`, ` `)
	inputText = s.Replace(inputText)
	return []byte(inputText), mapOfStringConstant
}

func replaceMultiSpaceToMonoSpace(s string) string {
	multiSpace := regexp.MustCompile(` +`)
	return multiSpace.ReplaceAllLiteralString(s, ` `)
}

func isSymbol(s string) bool {
	for _, sym := range symbols {
		if s == sym {
			return true
		}
	}
	return false
}

func isKeyword(s string) bool {
	for _, key := range keywords {
		if s == key {
			return true
		}
	}
	return false
}

func (tk *Tokenizer) GetCurrentToken() string {
	return tk.currentToken
}

func (tk *Tokenizer) GetTokens() []string {
	return tk.tokens
}
