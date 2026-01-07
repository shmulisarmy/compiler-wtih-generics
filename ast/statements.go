package ast

import (
	"fmt"

	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/display"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/runtime"
)

type Statement interface {
	Run()
}

type LastStatementAsExpression struct {
	Expression
}

func (this LastStatementAsExpression) Run() {
	display.DisplayStruct(this.Expression)
	panic("should hae used like an expression")
}

func (this If) Run() {
	this.Eval()
}
func (this FunctionCall) Run() {
	this.Eval()
}

func (this Block) Run() {
	this.Eval()
}

func (this Assignment) Run() {
	switch left := this.Left.(type) {
	case FieldAccess:
		object := left.Object.Eval()
		switch object := object.(type) {
		case map[string]runtime.LanguageValue:
			object[left.Name] = this.Right.Eval()
		default:
			panic("not implemented")
		}
	case Identifier:
		runtime.Set_value(left.Name, this.Right.Eval())
	default:
		panic("not implemented")
	}

}

type Assignment struct {
	Left  Expression
	Right Expression
	Range
}

func (i Assignment) String() string {
	return fmt.Sprintf("%s = %s",
		i.Left,
		i.Right)
}

func (this Return) Run() {
	panic("not implemented yet")
}
