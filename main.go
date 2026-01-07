package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/display"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/tokenizer"
)

func sanityCheck() {
	source := `:car.name`
	t := tokenizer.New(source)
	tokens := t.Tokenize()

	p := Parser{
		tokens: tokens,
	}
	for _, tok := range tokens {
		fmt.Println(tok)
	}
	e := p.Term()
	display.DisplayStruct(e)
	fmt.Println(e)
}

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
	e := p.Statement()
	fmt.Println(e)
	e.Run()
}
