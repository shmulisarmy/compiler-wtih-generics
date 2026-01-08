package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/assert"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/ast"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/tokenizer"
)

func main() {

	var source_file, err = os.Open("main.code")
	if err != nil {
		panic(err)
	}
	defer source_file.Close()

	source_bytes, err := ioutil.ReadAll(source_file)
	if err != nil {
		panic(err)
	}
	source := string(source_bytes)

	t := tokenizer.New(source)
	tokens := t.Tokenize()

	p := Parser{
		tokens: tokens,
	}
	for _, tok := range tokens {
		fmt.Println(tok)
	}

	globalScope := ast.GlobalScope{
		Symbols: map[string]ast.Symbol{},
	}

	for range 2 {
		p.Expect(tokenizer.VAR)
		var_ := p.GlobalVarDeclaration()
		globalScope.Symbols[var_.Name] = var_
	}
	block := p.Block()
	function := p.Function()
	globalScope.Symbols[function.Name] = function

	scopeStack := ast.ScopeStack{
		ast.ScopeFrame{Scope: &globalScope},
	}

	scopeStack.Type_check_block(&block)
	scopeStack.Type_check_function(&function)
	// scopeStack.Ensure_function_has_proper_return_type(&function)
	assert.Equal(len(scopeStack), 1)
	fmt.Println(function)
	fmt.Println(block)
	block.Run()
}
