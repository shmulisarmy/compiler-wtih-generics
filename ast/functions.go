package ast

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Function struct {
	Name        string
	Params      []VarDeclaration
	Body        Block
	Return_type string
	Range
}

func (this Function) String() string {
	var params []string
	for _, param := range this.Params {
		params = append(params, param.String())
	}
	return fmt.Sprintf("%s %s(%s) %s {\n%s\n%s\n%s",
		color.CyanString("Function"),
		this.Name,
		strings.Join(params, ", "),
		this.Return_type,
		this.Body,
		color.YellowString("}"),
		color.YellowString("}"),
	)
}
