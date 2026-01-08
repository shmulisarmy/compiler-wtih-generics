package ast

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/assert"
)

type ScopeFrame struct {
	scope      Scope
	stmt_index int
}
type ScopeStack []ScopeFrame

func (this ScopeStack) Type_check_function(function *Function) {
	this = append(this, ScopeFrame{scope: function, stmt_index: -1})
	for _, param := range function.Params {
		this.ensure_var_has_proper_type(&param)
	}
	this.Type_check_statements(function.Body.Statements)

}

func (this ScopeStack) Type_check_block(block *Block) {
	this = append(this, ScopeFrame{scope: block, stmt_index: -1})
	this.Type_check_statements(block.Statements)
}

func (this ScopeStack) ensure_var_has_proper_type(var_declaration *VarDeclaration) {
	if var_declaration.Type == "" {
		var_declaration.Type = get_expression_type(this, *var_declaration.DefaultValue.Expect(fmt.Sprintf("must have either a type or a default value to infer the type of %s", color.YellowString(var_declaration.Name))))
	} else if var_declaration.DefaultValue.IsPresent() {
		assert.Equal( //TODO:
			var_declaration.Type,
			get_expression_type(this, *var_declaration.DefaultValue.Unwrap()),
		)
	}
}

func (this ScopeStack) Symbol_name_to_type(symbol_name string) string {
	for i := len(this) - 1; i >= 0; i-- {
		sym := this[i].scope.Local_resolve(symbol_name, this[i].stmt_index)
		if sym != nil {
			switch sym := sym.(type) {
			case VarDeclaration:
				this.ensure_var_has_proper_type(&sym)
				return sym.Type
			}
		}
	}
	panic(fmt.Sprintf("couldn't find %s", color.YellowString(symbol_name)))

}

func (this ScopeStack) Type_check_statements(Statements []Statement) {
	for i, stmt := range Statements {
		switch stmt := stmt.(type) {
		case Block:
			previous_stmt_index := this[len(this)-1].stmt_index
			this[len(this)-1].stmt_index = i
			this.Type_check_block(&stmt)
			this[len(this)-1].stmt_index = previous_stmt_index
		case If:
			previous_stmt_index := this[len(this)-1].stmt_index
			this[len(this)-1].stmt_index = i
			this.Type_check_block(&stmt.Then)
			if stmt.Else.IsPresent() {
				this.Type_check_block(stmt.Else.Unwrap())
			}
			this[len(this)-1].stmt_index = previous_stmt_index
		case VarDeclaration:
			previous_stmt_index := this[len(this)-1].stmt_index
			this[len(this)-1].stmt_index = i
			this.ensure_var_has_proper_type(&stmt)
			this[len(this)-1].stmt_index = previous_stmt_index
		case Assignment:
			previous_stmt_index := this[len(this)-1].stmt_index
			this[len(this)-1].stmt_index = i
			left_type := get_expression_type(this, stmt.Left)
			right_type := get_expression_type(this, stmt.Right)
			this[len(this)-1].stmt_index = previous_stmt_index
			if left_type != right_type {
				panic(fmt.Sprintf("types do not match: %s at %s != %s at %s", color.GreenString(left_type), stmt.Left.PositionLink("main.code"), color.GreenString(right_type), stmt.Right.PositionLink("main.code")))
			}
		}
	}
}

// func (this *Block) Symbol_name_to_type(looking_for string) string {
// 	sym := this.Local_resolve(looking_for)
// 	switch sym := sym.(type) {
// 	case Assignment:
// 		return get_expression_type(this, sym.Right)
// 	default:
// 		panic("unhandled")
// 	}

// }
func (this *Block) Local_resolve(looking_for string, stmt_index int) Statement {
	for i, stmt := range this.Statements[:stmt_index] {
		print(i)
		switch stmt := stmt.(type) {
		case VarDeclaration:
			if stmt.Name == looking_for {
				return stmt
			}
		}
	}
	return nil
}

func (this *Block) resolve(var_name string, stmt_index int) Statement {
	stmt := this.Local_resolve(var_name, stmt_index)
	if stmt == nil {
		panic("coulnt find")
	}
	return stmt
}

type Scope interface {
	// resolve(var_name string, stmt_index int) Statement
	Local_resolve(looking_for string, parent_index int) Statement
	// Symbol_name_to_type(looking_for string) string
}

// /
func (this *Function) Local_resolve(var_name string, stmt_index int) Statement {
	stmt := this.Body.Local_resolve(var_name, stmt_index)
	if stmt != nil {
		return stmt
	}
	for _, param := range this.Params {
		if param.Name == var_name {
			return param
		}
	}
	return nil
}

///

func get_expression_type(scopes ScopeStack, e Expression) string {
	switch e := e.(type) {
	case Number:
		return "int"
	case String:
		return "string"
	case NoneValue:
		panic("not supposed to use this type as a value")
	case Block:
		switch last_stmt := e.Statements[len(e.Statements)-1].(type) {
		case LastStatementAsExpression:
			scopes = append(scopes, ScopeFrame{scope: &e, stmt_index: len(e.Statements) - 1})
			return get_expression_type(scopes, last_stmt.Expression)
		default:
			// return"none"
			panic("not supposed to use this type as a value")
		}
	case *Block:
		switch last_stmt := e.Statements[len(e.Statements)-1].(type) {
		case LastStatementAsExpression:
			scopes = append(scopes, ScopeFrame{scope: e, stmt_index: len(e.Statements) - 1})
			return get_expression_type(scopes, last_stmt.Expression)
		default:
			// return"none"
			panic("not supposed to use this type as a value")
		}
	case If:
		if e.Else.IsPresent() {
			then_type := get_expression_type(scopes, e.Then)
			else_type := get_expression_type(scopes, e.Else.Unwrap())
			if then_type != else_type {
				panic(fmt.Sprintf("types do not match: %s at %s != %s at %s", color.GreenString(then_type), e.Then.PositionLink("main.code"), color.GreenString(else_type), e.Else.Unwrap().PositionLink("main.code")))
			}
			return then_type
		} else {
			panic(fmt.Sprintf("in order to use an if statement as a value (%s) you must also have an else block", e.PositionLink("main.code")))
		}
		// panic("not supposed to use this type as a value")
	case ObjectInstanciation:
		return e.ClassName
	case BinaryExpression:
		left_type := get_expression_type(scopes, e.Left)
		right_type := get_expression_type(scopes, e.Right)
		if left_type != right_type {
			panic(fmt.Sprintf("types do not match: %s (at %s) and %s (at %s)", color.GreenString(left_type), e.Left.PositionLink("main.code"), color.GreenString(right_type), e.Right.PositionLink("main.code")))
		}
		return left_type
	case While:
		if e.ElseBlock.IsPresent() {
			then_type := get_expression_type(scopes, e.Body)
			else_type := get_expression_type(scopes, e.ElseBlock.Unwrap())
			if then_type != else_type {
				panic("types do not match")
			}
			return then_type
		}
		panic("not supposed to use this type as a value")
	case Identifier:
		return scopes.Symbol_name_to_type(e.Name)
	case FunctionCall:
		panic("cant yet use as type because there are no user defined functions")
	default:
		panic(fmt.Sprintf("unhandled: %s", e))
	}

}
