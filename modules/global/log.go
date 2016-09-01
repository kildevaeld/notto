package global

import (
	"fmt"
	"strings"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func getStringList(call otto.FunctionCall) []string {
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
			out = append(out, a.String())
		}

	}
	return out
}

func Define(vm *notto.Notto) error {
	if v, err := vm.Get("process"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	attr := vm.ProcessAttr()

	o, e := vm.Object("({})")
	if e != nil {
		return e
	}

	o.Set("log", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stdout != nil {
			attr.Stdout.Write([]byte(strings.Join(str, " ")))
		}

		return otto.UndefinedValue()
	})

	o.Set("error", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stderr != nil {
			attr.Stderr.Write([]byte(strings.Join(str, " ")))
		}

		return otto.UndefinedValue()
	})
}
