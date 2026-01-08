package ast

import "fmt"

type GlobalScope struct {
	Symbols map[string]Symbol
}

func (this GlobalScope) Local_resolve(looking_for string, stmt_index int) Symbol {
	if _, ok := this.Symbols[looking_for]; !ok {
		panic(fmt.Sprintf("cannot find %s", looking_for))
	}
	return this.Symbols[looking_for]
}
