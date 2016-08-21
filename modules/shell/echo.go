package shell

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func echoCall(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		if len(call.ArgumentList) == 0 {
			return otto.UndefinedValue()
		}

		var out []string
		for _, a := range call.ArgumentList {
			if a.IsString() {
				if s, e := a.ToString(); e == nil {
					out = append(out, s)
				}
			} else if a.IsNumber() {
				if s, e := a.ToInteger(); e == nil {
					out = append(out, fmt.Sprintf("%d", s))
				}
			} else {
				continue
			}

		}
		str := strings.Join(out, "")

		ob, _ := vm.Object("({})")

		fmt.Printf(str)

		ob.Set("pipe", func(call otto.FunctionCall) otto.Value {
			s, e := call.Argument(0).ToString()
			if e != nil {
				return otto.NullValue()
			}

			ioutil.WriteFile(s, []byte(str), 0755)

			return otto.UndefinedValue()
		})

		return ob.Value()
	}
}
