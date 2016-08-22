package util

import (
	"fmt"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func stringValue(str string) otto.Value {
	v, _ := otto.ToValue(str)
	return v
}

func formatCall(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			return stringValue("")
		}

		if !call.Argument(0).IsString() {
			vm.Throw("ArgumentError", "First arguments must be of type string")
		}

		str, err := call.Argument(0).ToString()
		if err != nil {
			panic(err)
		}

		var args []interface{}
		for _, a := range call.ArgumentList[1:] {
			v, e := a.Export()
			if e != nil {
				panic(e)
			}
			args = append(args, v)
		}

		out := fmt.Sprintf(str, args...)

		v, _ := vm.ToValue(out)
		return v
	}
}
