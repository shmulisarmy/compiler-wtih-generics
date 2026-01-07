package ast

import "github.com/shmuli/all-compilers/dk-what-to-call-this-one/runtime"

func (this FieldAccess) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	object := this.Object.Eval()
	switch object := object.(type) {
	case map[string]runtime.LanguageValue:
		return object[this.Name]
	default:
		panic("not implemented")
	}
}

func (this FunctionCall) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	function := this.FunctionReference.Eval()
	args := make([]runtime.LanguageValue, len(this.Args))
	for i, arg := range this.Args {
		args[i] = arg.Eval()
	}
	switch function := function.(type) {
	case func(...runtime.LanguageValue) runtime.LanguageValue:
		return function(args...)
	default:
		panic("not a built in function")
	}
}
func (this String) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	return this.Value
}
func (this Identifier) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	value := runtime.Get_value(this.Name)
	if value == nil {
		panic("not implemented")
	}
	return value
}
func (this Number) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	return this.Value
}
func (this BinaryExpression) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	switch this.Op {
	case "+":
		return this.Left.Eval().(int) + this.Right.Eval().(int)
	case "-":
		return this.Left.Eval().(int) - this.Right.Eval().(int)
	case "*":
		return this.Left.Eval().(int) * this.Right.Eval().(int)
	case "/":
		return this.Left.Eval().(int) / this.Right.Eval().(int)
	case "==":
		return this.Left.Eval() == this.Right.Eval()
	case "!=":
		return this.Left.Eval() != this.Right.Eval()
	case ">":
		return this.Left.Eval().(int) > this.Right.Eval().(int)
	case "<":
		return this.Left.Eval().(int) < this.Right.Eval().(int)
	case ">=":
		return this.Left.Eval().(int) >= this.Right.Eval().(int)
	case "<=":
		return this.Left.Eval().(int) <= this.Right.Eval().(int)
	default:
		panic("not implemented")
	}

}

func (this If) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	if this.Condition.Eval().(bool) {
		return this.Then.Eval()
	} else if this.Else.IsPresent() {
		elseBlock, _ := this.Else.Get()
		return elseBlock.Eval()
	}
	return NoneValue{}
}
func (this Block) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	for _, stmt := range this.Statements[:len(this.Statements)-1] {
		stmt.Run()
	}
	switch lastStmt := this.Statements[len(this.Statements)-1].(type) {
	case LastStatementAsExpression:
		return lastStmt.Eval()
	case Statement:
		lastStmt.Run()
		return NoneValue{}
	default:
		panic("not implemented")
	}

}
