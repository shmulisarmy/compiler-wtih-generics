package ast

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/assert"
)

type ScopeFrame struct {
	Scope      Scope
	stmt_index int
}
type ScopeStack []ScopeFrame

func (this VarDeclaration) Symbol__() {}
func (this ScopeStack) Ensure_function_has_proper_return_type(function *Function) {
	this = append(this, ScopeFrame{Scope: function, stmt_index: -1})
	for i, stmt := range function.Body.Statements {
		switch stmt := stmt.(type) {
		case Return:
			if !stmt.Value.IsPresent() {
				panic("return statement has no value")
			}
			this[len(this)-1].stmt_index = i
			if function.Return_type == "" {
				function.Return_type = get_expression_type(this, *stmt.Value.Unwrap())
			} else {
				assert.Equal(
					function.Return_type,
					get_expression_type(this, *stmt.Value.Unwrap()),
					fmt.Sprintf("return statement at %s has a value of type %s, but the %s function at %s has a return type of %s",
						stmt.PositionLink("main.code"),
						color.GreenString(get_expression_type(this, *stmt.Value.Unwrap())),
						color.YellowString(function.Name),
						function.PositionLink("main.code"),
						color.GreenString(function.Return_type)),
				)
			}
		}
	}
	assert.NotEqual(function.Return_type, "")
}

func (this ScopeStack) Type_check_function(function *Function) {
	l := len(this)
	this = append(this, ScopeFrame{Scope: function, stmt_index: -1})
	for _, param := range function.Params {
		this.ensure_var_has_proper_type(&param)
	}
	this.Type_check_statements(function.Body.Statements)
	this = this[:len(this)-1]
	assert.Equal(len(this), l)

}

func (this ScopeStack) Type_check_block(block *Block) {
	l := len(this)
	this = append(this, ScopeFrame{Scope: block, stmt_index: -1})
	this.Type_check_statements(block.Statements)
	this = this[:len(this)-1]
	assert.Equal(len(this), l)
}

func (this ScopeStack) ensure_var_has_proper_type(var_declaration *VarDeclaration) {
	// {
	// 	if var_declaration.Name == "global_num" {
	// 		assert.Equal(len(this), 1)
	// 	}
	// }
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
		sym := this[i].Scope.Local_resolve(symbol_name, this[i].stmt_index)
		if sym != nil {
			switch sym := sym.(type) {
			case VarDeclaration:
				this[:i+1].ensure_var_has_proper_type(&sym)
				return sym.Type
			case Function:
				this[:i+1].Ensure_function_has_proper_return_type(&sym)
				return "returns:" + sym.Return_type
			default:
				panic(fmt.Sprintf("unhandled: %s", sym))
			}
		}
	}
	panic(fmt.Sprintf("couldn't find %s", color.YellowString(symbol_name)))

}

func (this ScopeStack) Type_check_statements(Statements []Statement) {
	for i, stmt := range Statements {
		previous_stmt_index := this[len(this)-1].stmt_index
		assert.Equal(previous_stmt_index, -1)
		this[len(this)-1].stmt_index = i
		switch stmt := stmt.(type) {
		case Block:
			this.Type_check_block(&stmt)
		case If:
			this.Type_check_block(&stmt.Then)
			if stmt.Else.IsPresent() {
				this.Type_check_block(stmt.Else.Unwrap())
			}
		case VarDeclaration:
			this.ensure_var_has_proper_type(&stmt)
		case Assignment:
			left_type := get_expression_type(this, stmt.Left)
			right_type := get_expression_type(this, stmt.Right)
			if left_type != right_type {
				panic(fmt.Sprintf("types do not match: %s at %s != %s at %s", color.GreenString(left_type), stmt.Left.PositionLink("main.code"), color.GreenString(right_type), stmt.Right.PositionLink("main.code")))
			}
		case Return:
			//gets handled in the function typechecking
		case FunctionCall:
		default:
			panic(fmt.Sprintf("unhandled: %s", stmt))
		}
		this[len(this)-1].stmt_index = previous_stmt_index
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
func (this *Block) Local_resolve(looking_for string, stmt_index int) Symbol {
	if stmt_index == 4 {
		print(stmt_index)
	}
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

func (this *Block) resolve(var_name string, stmt_index int) Symbol {
	stmt := this.Local_resolve(var_name, stmt_index)
	if stmt == nil {
		panic("coulnt find")
	}
	return stmt
}

type Scope interface {
	// resolve(var_name string, stmt_index int) Statement
	Local_resolve(looking_for string, parent_index int) Symbol
	// Symbol_name_to_type(looking_for string) string
}

// /
func (this *Function) Local_resolve(var_name string, stmt_index int) Symbol {
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
			scopes = append(scopes, ScopeFrame{Scope: &e, stmt_index: len(e.Statements) - 1})
			t := get_expression_type(scopes, last_stmt.Expression)
			scopes = scopes[:len(scopes)-1]
			return t
		default:
			// return"none"
			panic("not supposed to use this type as a value")
		}
	case *Block:
		last_stmt := e.Statements[len(e.Statements)-1].(LastStatementAsExpression)
		l := len(scopes)
		scopes = append(scopes, ScopeFrame{Scope: e, stmt_index: len(e.Statements) - 1})
		t := get_expression_type(scopes, last_stmt.Expression)
		scopes = scopes[:len(scopes)-1]
		assert.Equal(len(scopes), l)
		return t
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
		not_a_template := true
		assert.Assert(not_a_template)
		function := scopes.Symbol_name_to_type(e.FunctionReference.(Identifier).Name)
		assert.Assert(strings.HasPrefix(function, "returns:"))
		return function[len("returns:"):]
	default:
		panic(fmt.Sprintf("unhandled: %s", e))
	}

}
