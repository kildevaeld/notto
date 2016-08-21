package shell

import (
	"io/ioutil"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func catCall(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		var str string
		var err error
		if str, err = call.Argument(0).ToString(); err != nil {
			return otto.NullValue()
		}

		bs, e := ioutil.ReadFile(str)
		if e != nil {
			return otto.NullValue()
		}

		v, _ := otto.ToValue(string(bs))

		return v
	}
}
