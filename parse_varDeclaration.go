package main

import (
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/ast"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/optional"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/tokenizer"
)

func (this *Parser) GlobalVarDeclaration() VarDeclaration {
	start_pos := this.curPos()
	var_name := this.Expect(IDENT).Literal
	var type_ string
	if this.currentToken().Type == IDENT {
		type_ = this.Expect(IDENT).Literal
	}
	default_value := None[Expression]()
	if this.OptionallyExpect(ASSIGN) {
		default_value = Some(this.Expression())
	}
	return VarDeclaration{
		Name:         var_name,
		Type:         type_,
		DefaultValue: default_value,
		Range: Range{
			Start: start_pos,
			End:   this.curPos(),
		},
	}
}
