rewriting-JackCompiler: compilationengine/*.go jacktokenizer/*.go symboltable/*.go vmwriter/*.go main.go
	go build
clean:
	rm -f rewriting-JackCompiler
	rm -f testcases/*/*.vm
	rm -f testcases/*/*.xml
test: rewriting-JackCompiler
	bash test.sh
