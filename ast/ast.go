package ast

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/optional"
	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/runtime"
)

// ============================================================================
// Position and Range Types
// ============================================================================

type Pos struct {
	Line int
	Col  int
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%s%d%s:%s%d%s",
		color.CyanString("Pos"),
		color.YellowString(""),
		p.Line,
		color.YellowString(""),
		color.YellowString(""),
		p.Col,
		color.YellowString(""))
}

type Range struct {
	Start Pos
	End   Pos
}

func (r Range) PositionLink(filename string) string {
	return color.BlueString(fmt.Sprintf("%s:%d:%d", filename, r.Start.Line, r.Start.Col))
}

func (r Range) String() string {
	return fmt.Sprintf("%s[%s → %s]",
		color.CyanString("Range"),
		color.GreenString("%d:%d", r.Start.Line, r.Start.Col),
		color.GreenString("%d:%d", r.End.Line, r.End.Col))
}

// ============================================================================
// Expression Interface and Implementations
// ============================================================================

type Expression interface {
	Eval() runtime.LanguageValue
	PositionLink(filename string) string
}

// Expression implementations
func (Number) Expression__()           {}
func (UnaryExpression) Expression__()  {}
func (BinaryExpression) Expression__() {}
func (Identifier) Expression__()       {}

type Identifier struct {
	Name string
	Range
}

func (i Identifier) String() string {
	return fmt.Sprintf("%s(%s)",
		color.CyanString("Identifier"),
		color.GreenString(i.Name))
}

type Number struct {
	Value int
	Range
}

func (l Number) String() string {
	return fmt.Sprintf("%d", l.Value)
}

type String struct {
	Value string
	Range
}

func (l String) String() string {
	return fmt.Sprintf("%s", l.Value)
}

type UnaryExpression struct {
	Op   string
	Expr Expression
}

func (u UnaryExpression) String() string {
	return fmt.Sprintf("%s(%s %v)",
		color.CyanString("UnaryExpr"),
		color.MagentaString(u.Op),
		u.Expr)
}

type BinaryExpression struct {
	Left  Expression
	Op    string
	Right Expression
	Range
}

func (b BinaryExpression) String() string {
	return fmt.Sprintf("(%v %s %v)",
		b.Left,
		color.MagentaString(b.Op),
		b.Right)
}

// ============================================================================
// Statement Interface and Implementations
// ============================================================================

func (Block) Expression__() {}
func (Block) Statement__()  {}

type Symbol interface{ Symbol__() }

func (Assignment) Symbol__() {}

type Block struct {
	Statements []Statement
	Range
}

func (b Block) String() string {
	var stmts []string
	for _, stmt := range b.Statements {
		stmts = append(stmts, fmt.Sprintf("  %v", stmt))
	}
	stmtStr := strings.Join(stmts, "\n")
	if stmtStr == "" {
		return color.CyanString("Block") + color.YellowString(" {}")
	}
	return fmt.Sprintf("%s %s\n%s\n%s",
		color.CyanString("Block"),
		color.YellowString("{"),
		stmtStr,
		color.YellowString("}"))
}

// ============================================================================
// Other Types
// ============================================================================

func (FunctionCall) Expression__() {

}

type FunctionCall struct {
	FunctionReference Expression
	Args              []Expression
	CatchBlock        optional.Optional[Block]
	Range
}

func (f FunctionCall) String() string {
	args := ""
	for _, arg := range f.Args {
		args += fmt.Sprintf("%s, ", arg)
	}
	if len(args) > 2 {
		args = args[:len(args)-2]
	}
	result := fmt.Sprintf("%s %s(%s)",
		color.CyanString("FunctionCall"),
		f.FunctionReference,
		color.GreenString(args))

	if f.CatchBlock.IsPresent() {
		block, _ := f.CatchBlock.Get()
		result += fmt.Sprintf("\n  %s %v",

			color.YellowString("catch"),
			block)
	}

	return result
}

func (FieldAccess) Expression__() {

}

type FieldAccess struct {
	Object Expression
	Name   string
	Range
}

func (f FieldAccess) String() string {
	return fmt.Sprintf("%s.%s",
		f.Object,
		color.BlueString(f.Name))
}

type Return struct {
	Value optional.Optional[Expression]
	Range
}

func (r Return) String() string {
	return fmt.Sprintf("%s %v",
		color.CyanString("Return"),
		r.Value)
}

type If struct {
	Condition Expression
	Then      Block
	Else      optional.Optional[Block]
	Range
}

func (If) Expression__() {

}
func (i If) String() string {
	var elseStr string
	if i.Else.IsPresent() {
		elseBlock, _ := i.Else.Get()
		elseStr = fmt.Sprintf("\n  %s %v",
			color.YellowString("else"),
			elseBlock)
	}
	return fmt.Sprintf("%s %v\n%s\n%s",
		color.CyanString("If"),
		i.Condition,
		i.Then,
		elseStr)
}

type While struct {
	Condition Expression
	Body      Block
	ElseBlock optional.Optional[Block]
	Range
}

func (w While) String() string {
	return fmt.Sprintf("%s %v\n%s",
		color.CyanString("While"),
		w.Condition,
		w.Body)
}

func (this String) Expression__() {
}

type NoneValue struct {
	Range
}

func (NoneValue) String() string {
	return color.CyanString("None")
}

func (NoneValue) Expression__() {

}
func (NoneValue) Eval() runtime.LanguageValue {
	panic("this is not a proper value")
}

type FieldNode struct {
	Name  string
	Value Expression
	Range
}
type ObjectInstanciation struct {
	ClassName string
	Fields    []FieldNode
	Range
}

func (this ObjectInstanciation) String() string {
	var fields []string
	for _, field := range this.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", field.Name, field.Value))
	}
	return fmt.Sprintf("%s{%s}", color.GreenString(this.ClassName),
		strings.Join(fields, ", "))
}
func (this ObjectInstanciation) Eval() runtime.LanguageValue {
	object := make(map[string]runtime.LanguageValue)
	for _, field := range this.Fields {
		object[field.Name] = field.Value.Eval()
	}
	return object
}

func (this ObjectInstanciation) Expression__() {

}

type Using struct {
	Objects Expression
	Block   Block
	Range
}

func (this Using) String() string {
	return fmt.Sprintf("%s %v\n%s",
		color.CyanString("Using"),
		this.Objects,
		this.Block)
}

func (this Using) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	object := this.Objects.Eval()
	switch object := object.(type) {
	case map[string]runtime.LanguageValue:
		runtime.UsingStack = append(runtime.UsingStack, object)
		defer func() {
			runtime.UsingStack = runtime.UsingStack[:len(runtime.UsingStack)-1]
		}()
		return this.Block.Eval()
	default:
		panic("not implemented")
	}
}

func (this Using) Expression__() {

}
func (this Using) Run() {
	this.Eval()
}

type UsingField struct {
	Name string
	Range
}

func (this UsingField) String() string {
	return fmt.Sprintf("%s.%s",
		this.Name,
		color.BlueString(this.Name))
}

func (this UsingField) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	for i := len(runtime.UsingStack) - 1; i >= 0; i-- {
		if object, ok := runtime.UsingStack[i][this.Name]; ok {
			return object
		}
	}
	panic("not implemented")
}

func (this While) Eval() runtime.LanguageValue {
	print(HighlightRange("", this.Range))
	for this.Condition.Eval().(bool) {
		for _, stmt := range this.Body.Statements {

			stmt.Run()
			if runtime.ControlFlowContext.ShouldBreak {
				runtime.ControlFlowContext.ShouldBreak = false
				if runtime.ControlFlowContext.BreakValue.IsPresent() {
					v := runtime.ControlFlowContext.BreakValue.Unwrap()
					runtime.ControlFlowContext.BreakValue = optional.None[runtime.LanguageValue]()
					return v
				}
				return NoneValue{}
			}
		}
	}
	if this.ElseBlock.IsPresent() {
		elseBlock, _ := this.ElseBlock.Get()
		v := elseBlock.Eval()
		if (v == NoneValue{}) {
			panic(fmt.Sprintf(""))

		}
		return v
	}
	return NoneValue{}
}

///

type Instruction struct {
	Operation func()
	Range
}

func (this Instruction) String() string {
	return fmt.Sprintf("%s",
		color.CyanString("Instruction"),
		this.Range)
}

func (this Instruction) Run() {
	this.Operation()
}

type Unless struct {
	Condition Expression
	Body      Statement
	Range
}

func (this Unless) String() string {
	return fmt.Sprintf("%s %v\n%s",
		color.CyanString("Unless"),
		this.Condition,
		this.Body)
}

func (this Unless) Run() {
	if this.Condition.Eval().(bool) {
		return
	}
	this.Body.Run()
}

//

type VarDeclaration struct {
	Name         string
	Type         string
	DefaultValue optional.Optional[Expression]
	Range
}

func (this VarDeclaration) String() string {
	var default_value string
	if this.DefaultValue.IsPresent() {
		default_value = fmt.Sprintf(" = %s", *this.DefaultValue.Unwrap())
	}
	return fmt.Sprintf("%s%s%s",
		color.YellowString(this.Name),
		color.GreenString(this.Type),
		default_value)
}

func (this VarDeclaration) Run() {
	if this.DefaultValue.IsPresent() {
		runtime.Set_value(this.Name, (*this.DefaultValue.Unwrap()).Eval())
	}

}
