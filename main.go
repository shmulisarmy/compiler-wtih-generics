package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	block := p.Block()
	function := p.Function()

	ast.ScopeStack{}.Type_check_block(&block)
	ast.ScopeStack{}.Type_check_function(&function)
	fmt.Println(function)
	fmt.Println(block)
	block.Run()
}
