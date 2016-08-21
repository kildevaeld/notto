package shell

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

func pipeFn(str string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		if call.Argument(0).IsFunction() {

			vv, _ := otto.ToValue(str)
			v, _ := call.Argument(0).Call(otto.NullValue(), vv)

			return v
		}
		s, e := call.Argument(0).ToString()
		if e != nil {
			return otto.NullValue()
		}

		ioutil.WriteFile(s, []byte(str), 0755)

		return otto.UndefinedValue()
	}
}
