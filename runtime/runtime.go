package runtime

import (
	"fmt"

	"github.com/shmuli/all-compilers/dk-what-to-call-this-one/optional"
)

type LanguageValue any

var Vars = map[string]LanguageValue{
	"true":  true,
	"false": false,
	"age":   25,
	"print": func(args ...LanguageValue) LanguageValue {
		for i, arg := range args {
			fmt.Print(arg)
			if i != len(args)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		return nil
	},
	"output": func(args ...LanguageValue) LanguageValue {
		for i, arg := range args {
			fmt.Print(arg)
			if i != len(args)-1 {
				fmt.Print(" ")
			}
		}
		return nil
	},
}

var UsingStack []map[string]LanguageValue

func Get_value(name string) LanguageValue {
	return Vars[name]
}

func Set_value(name string, value LanguageValue) {
	Vars[name] = value
}

type ControlFlowContextType struct {
	ShouldContinue bool
	ShouldBreak    bool
	ShouldReturn   bool
	ReturnValue    optional.Optional[LanguageValue]
	BreakValue     optional.Optional[LanguageValue]
}

var ControlFlowContext = ControlFlowContextType{
	ShouldContinue: false,
	ShouldBreak:    false,
	ShouldReturn:   false,
	ReturnValue:    optional.None[LanguageValue](),
	BreakValue:     optional.None[LanguageValue](),
}
