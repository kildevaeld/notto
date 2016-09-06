//go:generate go-bindata -pkg ui -o ui_impl.go ui.js
package ui

import (
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/prompt"
	"github.com/robertkrimen/otto"
)

type Map map[string]otto.Value

func mustValue(v otto.Value, err error) otto.Value {
	if err != nil {
		panic(err)
	}
	return v
}

func Define(vm *notto.Notto) error {

	var (
		o   *otto.Object
		err error
	)

	if o, err = vm.Object("({})"); err != nil {
		return err
	}

	o.Set("form", form(vm))
	o.Set("process", process(vm))
	loader := notto.CreateLoaderFromSourceWithPrivate(string(MustAsset("ui.js")), "", Map{
		"__private_ui": o.Value(),
	})

	vm.AddModule("ui", loader)

	return nil
}

func form(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		o := call.Argument(0).Object()
		cui := prompt.NewUI()
		t := mustValue(o.Get("type")).String()
		msg := mustValue(o.Get("message")).String()

		var val interface{}
		switch t {
		case "input":
			val = cui.Input(msg)
		case "password":
			val = cui.Password(msg)

		default:
			vm.Throw("UIError", "input type does not exists")
		}

		if v, e := vm.ToValue(val); e != nil {
			vm.Throw("UIError", e)
		} else {
			return v
		}

		return otto.UndefinedValue()
	}
}

func process(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		msg := call.Argument(0).String()
		cui := prompt.NewUI()

		p := cui.Process(msg)

		if v, e := vm.ToValue(p); e != nil {
			vm.Throw("UIError", e)
		} else {
			return v
		}
		return otto.UndefinedValue()
	}
}
