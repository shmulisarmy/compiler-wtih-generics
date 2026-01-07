package main

import (
	"fmt"
	"strconv"

	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/ast"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/optional"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/runtime"
	. "github.com/shmuli/all-compilers/dk-what-to-call-this-one/tokenizer"
)

type Parser struct {
	tokens []Token
	pos    int
}

func (this *Parser) Expect(expectedType TokenType) Token {
	t := this.TakeToken()
	if t.Type != expectedType {
		panic(fmt.Sprintf("expected %s, got %s", expectedType, t.Type))
	}
	return t
}
func (this *Parser) TakeToken() Token {
	tok := this.tokens[this.pos]
	this.pos++
	return tok
}

func (this *Parser) OptionallyExpect(expectedType TokenType) bool {
	t := this.currentToken()
	if t.Type == expectedType {
		this.TakeToken()
		return true
	}
	return false

}
func (this *Parser) currentToken() Token {
	return this.tokens[this.pos]
}

func (this *Parser) curPos() Pos {
	return Pos{Line: this.currentToken().Line, Col: this.currentToken().Column}
}
func (this *Parser) Statement() Statement {
	start_pos := this.curPos()
	if this.OptionallyExpect(RETURN) {
		return Return{
			Value: Some(this.Expression()),
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
	}
	if this.OptionallyExpect(BREAK) {
		I := Instruction{
			Operation: func() {
				runtime.ControlFlowContext.ShouldBreak = true
			},
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
		if this.OptionallyExpect(COLON) {
			e := this.Expression()
			I.Operation = func() {
				runtime.ControlFlowContext.ShouldBreak = true
				runtime.ControlFlowContext.BreakValue = Some(e.Eval())
			}
			I.Range.End = this.curPos()
		}
		return I
	}
	e := this.Expression()
	switch e := e.(type) {
	case Block:
		return e
	case If:
		return e
	case Using:
		return e
	case FunctionCall:
		return e
	case FieldAccess, Identifier:
		if this.OptionallyExpect(ASSIGN) {
			return Assignment{
				Left:  e,
				Right: this.Expression(),
				Range: Range{
					Start: start_pos,
					End:   this.curPos(),
				},
			}
		}
	}
	return LastStatementAsExpression{Expression: e}
}
func (this *Parser) Block() Block {
	var stmts []Statement
	this.Expect(LBRACE)
	for !this.OptionallyExpect(RBRACE) {
		stmt := this.Statement()
		if this.OptionallyExpect(UNLESS) {
			start_pos := this.curPos()
			stmt = Unless{
				Condition: this.Expression(),
				Body:      stmt,
				Range: Range{
					Start: start_pos,
					End:   this.curPos(),
				},
			}
		}
		stmts = append(stmts, stmt)
	}
	return Block{
		Statements: stmts,
		Range: Range{
			Start: this.curPos(),
			End:   this.curPos(),
		},
	}
}

func (this *Parser) Term() Expression {

	start_pos := this.curPos()
	if this.currentToken().Type == LBRACE {
		fields := CustomList(this, LBRACE, RBRACE, func(this *Parser) FieldNode {
			field_name := this.Expect(IDENT).Literal
			this.Expect(ASSIGN)
			return FieldNode{
				Name:  field_name,
				Value: this.Expression(),
				Range: Range{
					Start: this.curPos(),
					End:   this.curPos(),
				},
			}
		})
		return ObjectInstanciation(fields)
	}
	t := this.TakeToken()
	switch t.Type {
	case USING:
		return Using{
			Objects: this.Expression(),
			Block:   this.Block(),
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
	case WHILE:
		w := While{
			Condition: this.Expression(),
			Body:      this.Block(),
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
		if this.OptionallyExpect(ELSE) {
			w.ElseBlock = Some(this.Block())
		}
		return w
	case IF:
		if_ := If{
			Condition: this.Expression(),
			Then:      this.Block(),
			Else:      None[Block](),
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}
		if this.OptionallyExpect(ELSE) {
			if_.Else = Some(this.Block())
		}
		return if_
	case LPAREN:
		e := this.Expression()
		this.Expect(RPAREN)
		return e

	case IDENT, COLON:
		var expr Expression
		switch t.Type {
		case IDENT:
			expr = Identifier{
				Name: t.Literal,
				Range: Range{
					Start: start_pos,
					End:   this.curPos(),
				},
			}
		case COLON:
			expr = UsingField{
				Name: this.Expect(IDENT).Literal,
				Range: Range{
					Start: start_pos,
					End:   this.curPos(),
				},
			}
		default:
			panic(fmt.Sprintf("not implemented: %s", t))
		}
		for {
			t := this.currentToken()
			if t.Type == DOT {
				this.TakeToken()
				t = this.Expect(IDENT)
				expr = FieldAccess{
					Object: expr,
					Name:   t.Literal,
					Range: Range{
						Start: start_pos,
						End:   this.curPos(),
					},
				}
			} else if t.Type == LPAREN {
				return FunctionCall{
					FunctionReference: expr,
					Args:              this.ExpressionList(LPAREN, RPAREN),
					CatchBlock:        None[Block](),
					Range: Range{
						Start: start_pos,
						End:   this.curPos(),
					},
				}
			} else {
				break
			}
		}
		return expr

	case INT:
		n, _ := strconv.Atoi(t.Literal)
		return Number{
			Value: n,
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}

	case STRING:
		return String{
			Value: t.Literal,
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}

	default:
		panic(fmt.Sprintf("not implemented: %s", t))
	}
	panic(fmt.Sprintf("not implemented: %s", t))

}

var operatorPrecedence = map[TokenType]int{
	EQ:    1,
	NEQ:   1,
	GT:    1,
	LT:    1,
	GTE:   1,
	LTE:   1,
	AND:   1,
	OR:    1,
	PLUS:  2,
	MINUS: 2,
	MUL:   3,
	DIV:   3,
}

func CustomList[T any](this *Parser, openingToken TokenType, closingToken TokenType, f func(this *Parser) T) []T {
	exprs := []T{}
	this.Expect(openingToken)
	for !this.OptionallyExpect(closingToken) {
		exprs = append(exprs, f(this))
		if !this.OptionallyExpect(COMMA) {
			this.Expect(closingToken)
			break
		}
	}

	return exprs
}
func (this *Parser) ExpressionList(openingToken TokenType, closingToken TokenType) []Expression {
	exprs := []Expression{}
	this.Expect(openingToken)
	for !this.OptionallyExpect(closingToken) {
		expr := this.Expression()
		exprs = append(exprs, expr)
		if !this.OptionallyExpect(COMMA) {
			this.Expect(closingToken)
			break
		}
	}

	return exprs
}
func (this *Parser) Expression() Expression {
	start_pos := this.curPos()
	left := this.Term()
	for {
		t := this.currentToken()
		if _, ok := operatorPrecedence[t.Type]; !ok {
			break
		}
		this.TakeToken()
		right := this.Expression()
		left = BinaryExpression{
			Left:  left,
			Op:    t.Literal,
			Right: right,
			Range: Range{
				Start: start_pos,
				End:   this.curPos(),
			},
		}

	}
	return left

}
