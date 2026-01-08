package main

import (
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/ast"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/optional"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/tokenizer"
)

func (this *Parser) Function() Function {
	start_pos := this.curPos()
	this.Expect(FUNC)
	name := this.Expect(IDENT).Literal
	params := ParseCustomList(this, LPAREN, RPAREN, func(this *Parser) VarDeclaration {
		start_pos := this.curPos()
		return VarDeclaration{
			Name:         this.Expect(IDENT).Literal,
			Type:         this.Expect(IDENT).Literal,
			DefaultValue: None[Expression](),
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
	})
	return_type := this.Expect(IDENT).Literal
	return Function{
		Name:        name,
		Params:      params,
		Return_type: return_type,
		Body:        this.Block(),
		Range: Range{
			Start: start_pos,
			End:   this.curPos(),
		},
	}
}
