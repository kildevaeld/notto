//go:generate go-bindata -pkg crypto -o crypto_impl.go crypto.js
package crypto

import (
	"crypto/rand"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/kildevaeld/notto/modules/global"
	"github.com/robertkrimen/otto"
)

type simple_task struct {
	id     int64
	err    error
	result []byte
	call   otto.Value
	vm     *notto.Notto
}

func (self *simple_task) SetID(id int64) { self.id = id }
func (self *simple_task) GetID() int64   { return self.id }

func (self *simple_task) Execute(vm *otto.Otto, loop *loop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		/*e, err := vm.Call(`new Error`, nil, self.err.Error())
		if err != nil {
			return err
		}*/
		e := vm.MakeCustomError("CryptoError", self.err.Error())

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}

	/*if self.result != nil {
		if v, e := vm.ToValue(t.content); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}*/
	arguments = append(arguments, global.NewBuffer(self.vm, self.result))

	//arguments = append([]interface{}{self.call}, arguments...)
	if _, err := self.call.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}
	/*if _, err := vm.Call(`Function.call.call`, nil, arguments...); err != nil {
		return err
	}**/

	return nil
}

func (self *simple_task) Cancel() {

}

func Define(vm *notto.Notto) error {

	o, _ := vm.Object("({})")
	o.Set("randomBytes", func(call otto.FunctionCall) otto.Value {

		l, e := call.Argument(0).ToInteger()
		if e != nil {
			vm.Throw("CryptoError", "length")
		}

		task := &simple_task{vm: vm}
		task.call = call.Argument(1)
		vm.Runloop().Add(task)
		go func() {
			defer vm.Runloop().Ready(task)

			b := make([]byte, l)

			_, e := rand.Read(b)
			if e != nil {
				task.err = e
			} else {
				//v, _ := vm.ToValue(string(b))
				task.result = b // v //global.NewBuffer(vm, b)
			}

		}()
		/*looptask.SimpleTask(vm, call.Argument(1), func(t loop.Task) (interface{}, error) {
			b := make([]byte, l)

			_, e := rand.Read(b)
			if e != nil {
				return nil, e
			}

			return global.NewBuffer(vm, b), nil
		})*/

		return otto.UndefinedValue()
	})

	m := make(map[string]otto.Value)
	m["__private_crypto"] = o.Value()

	vm.AddModule("crypto", notto.CreateLoaderFromSourceWithPrivate(string(MustAsset("crypto.js")), "", m))

	return nil
}
