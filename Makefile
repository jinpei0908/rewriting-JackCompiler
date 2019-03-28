JackCompiler: compilationengine/*.go jacktokenizer/*.go symboltable/*.go vmwriter/*.go main.go
	go build
clean:
	rm -f JackCompiler 
