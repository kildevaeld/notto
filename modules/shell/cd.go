package shell

import (
	"os"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func cdCall(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		s, e := call.Argument(0).ToString()
		if e != nil {
			return otto.UndefinedValue()
		}
		os.Chdir(s)

		return otto.UndefinedValue()
	}

}
