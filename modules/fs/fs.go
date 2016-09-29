//go:generate go-bindata -pkg fs -o fs_impl.go fs.js
package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

type fs_task struct {
	id      int64
	cb      otto.Value
	err     error
	content interface{}
}

func (t *fs_task) SetID(id int64) { t.id = id }
func (t *fs_task) GetID() int64   { return t.id }

func (t *fs_task) Execute(vm *otto.Otto, loop *loop.Loop) error {
	var arguments []interface{}

	if t.err != nil {
		e, err := vm.Call(`new Error`, nil, t.err.Error())
		if err != nil {
			return err
		}

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}
	if t.content != nil {
		if v, e := vm.ToValue(t.content); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}

	if _, err := t.cb.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}

	return nil
}

func (self *fs_task) Cancel() {

}

func _set(result error) func(err error) error {
	return func(err error) error {
		if err != nil {
			return multierror.Append(result, err)
		}
		return nil
	}
}

func Define(vm *notto.Notto) error {

	var (
		o   *otto.Object
		err error
	)

	if o, err = vm.Object("({})"); err != nil {
		return err
	}

	var result error

	err = _set(result)(o.Set("readFile", readFile(vm)))
	err = _set(result)(o.Set("readdir", readdir(vm)))
	err = _set(result)(o.Set("rename", rename(vm)))
	err = _set(result)(o.Set("stat", stat(vm)))
	if err != nil {
		return err
	}

	vm.AddModule("__private_fs", notto.CreateLoaderFromValue(o.Value()))
	vm.AddModule("fs", notto.CreateLoaderFromSource(string(MustAsset("fs.js")), ""))
	return nil
}

func stat(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		//cwd := vm.ProcessAttr().Cwd
		source := call.Argument(0).String()

		t := &fs_task{
			cb: call.Argument(1),
		}

		vm.Runloop().Add(t)
		go func() {
			defer vm.Runloop().Ready(t)

			bs, err := os.Stat(source)
			if err != nil {
				t.err = err
			} else {
				t.content = bs
			}
		}()

		return otto.UndefinedValue()
	}
}

func readFile(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		cwd := vm.ProcessAttr().Cwd
		source := call.Argument(0).String()
		source = filepath.Join(cwd, source)

		t := &fs_task{
			cb: call.Argument(1),
		}

		vm.Runloop().Add(t)
		go func() {
			defer vm.Runloop().Ready(t)

			bs, err := ioutil.ReadFile(source)
			if err != nil {
				t.err = err
			} else {
				t.content = string(bs)
			}
		}()

		return otto.UndefinedValue()
	}
}

func readdir(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		cwd := vm.ProcessAttr().Cwd
		source := call.Argument(0).String()
		source = filepath.Join(cwd, source)

		t := &fs_task{
			cb: call.Argument(1),
		}

		vm.Runloop().Add(t)
		go func() {
			defer vm.Runloop().Ready(t)

			files, err := ioutil.ReadDir(source)
			if err != nil {
				t.err = err
			} else {
				t.content = files
			}
		}()

		return otto.UndefinedValue()

	}
}

func rename(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		cwd := vm.ProcessAttr().Cwd
		source := call.Argument(0).String()
		source = filepath.Join(cwd, source)
		target := filepath.Join(cwd, call.Argument(1).String())

		t := &fs_task{
			cb: call.Argument(2),
		}

		vm.Runloop().Add(t)
		go func() {
			defer vm.Runloop().Ready(t)
			t.err = os.Rename(source, target)
		}()

		return otto.UndefinedValue()

	}
}
