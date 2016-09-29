package process

import (
	"os"
	"runtime"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func stringToValue(str string) otto.Value {
	v, e := otto.ToValue(str)
	if e != nil {
		panic(e)
	}
	return v
}

func DefineFunc(vm *notto.Notto, fn func() (notto.Environ, []string)) error {
	if v, err := vm.Get("__private_process"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	vm.AddPreScript(`if (typeof process === 'undefined') var process;
		process = __private_process();
		`)

	attr := vm.ProcessAttr()

	return vm.Set("__private_process", func(call otto.FunctionCall) otto.Value {

		env, argv := fn()

		v, _ := vm.ToValue(map[string]interface{}{
			"env":  env,
			"argv": argv,
			"cwd": func(call otto.FunctionCall) otto.Value {
				return stringToValue(attr.Cwd)
			},
			"platform": runtime.GOOS,
			"exit": func(call otto.FunctionCall) {
				os.Exit(0)
			},
		})

		return v

	})

}

func Define(vm *notto.Notto) error {
	if v, err := vm.Get("process"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	attr := vm.ProcessAttr()

	return vm.Set("process", map[string]interface{}{
		"env":  attr.Environ.ToMap(),
		"argv": attr.Argv,
		"cwd": func(call otto.FunctionCall) otto.Value {
			return stringToValue(attr.Cwd)
		},
		"platform": runtime.GOOS,
		"exit": func(n int) {
			os.Exit(n)
		},
	})
}
