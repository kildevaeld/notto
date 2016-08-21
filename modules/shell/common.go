package shell

import "github.com/robertkrimen/otto"

func pipeFn(str string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		if call.Argument(0).IsFunction() {

			vv, _ := vm.ToValue(str)
			v, _ := call.Argument(0).Call(otto.NullValue(), vv)

			return v
		}
		s, e := call.Argument(0).ToString()
		if e != nil {
			return otto.NullValue()
		}

		ioutil.WriteFile(s, []byte(r.stdout), 0755)

		return otto.UndefinedValue()
	}
}
