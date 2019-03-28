package symboltable

import "errors"

var err error = errors.New("Error occured. ")

type valueOfHashMap struct {
	type_ string
	kind  string
	index int
}

type SymbolTable struct {
	TableOfClassScope      map[string]valueOfHashMap
	TableOfSubroutineScope map[string]valueOfHashMap
	CurrentKind            string
	CurrentType            string
	staticIndex            int
	fieldIndex             int
	varIndex               int
	argIndex               int
	IfCount                int
	WhileCount             int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		TableOfClassScope:      map[string]valueOfHashMap{},
		TableOfSubroutineScope: map[string]valueOfHashMap{},
		staticIndex:            0,
		fieldIndex:             0,
		varIndex:               0,
		argIndex:               0,
		IfCount:                0,
		WhileCount:             0}
}

func (st *SymbolTable) StartSubroutine(subroutineKind string) {
	st.TableOfSubroutineScope = map[string]valueOfHashMap{}
	st.varIndex = 0
	st.argIndex = 0
	st.IfCount = 0
	st.WhileCount = 0
	if subroutineKind == "method" {
		st.Define("this", "", "arg")
	}
}

func (st *SymbolTable) Define(name, type_, kind string) {
	switch kind {
	case "static":
		st.TableOfClassScope[name] = newValue(type_, "Static", st.staticIndex)
		st.staticIndex++
	case "field":
		st.TableOfClassScope[name] = newValue(type_, "Field", st.fieldIndex)
		st.fieldIndex++
	case "var":
		st.TableOfSubroutineScope[name] = newValue(type_, "Var", st.varIndex)
		st.varIndex++
	case "arg":
		st.TableOfSubroutineScope[name] = newValue(type_, "Argument", st.argIndex)
		st.argIndex++
	}
}

func (st *SymbolTable) VarCount(kind string) (int, error) {
	switch kind {
	case "Static":
		return st.staticIndex, nil
	case "Field":
		return st.fieldIndex, nil
	case "Var":
		return st.varIndex, nil
	case "Argument":
		return st.argIndex, nil
	default:
		return -1, err
	}
}

func (st *SymbolTable) KindOf(name string) (string, error) {
	valueInSubroutineScope, ok := st.TableOfSubroutineScope[name]
	if ok {
		return valueInSubroutineScope.kind, nil
	}
	valueInClassScope, ok := st.TableOfClassScope[name]
	if ok {
		return valueInClassScope.kind, nil
	}
	return "", err
}

func (st *SymbolTable) TypeOf(name string) (string, error) {
	valueInSubroutineScope, ok := st.TableOfSubroutineScope[name]
	if ok {
		return valueInSubroutineScope.type_, nil
	}
	valueInClassScope, ok := st.TableOfClassScope[name]
	if ok {
		return valueInClassScope.type_, nil
	}
	return "", err
}

func (st *SymbolTable) IndexOf(name string) (int, error) {
	valueInSubroutineScope, ok := st.TableOfSubroutineScope[name]
	if ok {
		return valueInSubroutineScope.index, nil
	}
	valueInClassScope, ok := st.TableOfClassScope[name]
	if ok {
		return valueInClassScope.index, nil
	}
	return -1, err
}

func newValue(type_, kind string, index int) valueOfHashMap {
	return valueOfHashMap{
		type_: type_,
		kind:  kind,
		index: index}
}
