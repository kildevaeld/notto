package looptask

import (
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

type simple_task struct {
	id     int64
	err    error
	result otto.Value
	call   otto.Value
}

func (self *simple_task) SetID(id int64) { self.id = id }
func (self *simple_task) GetID() int64   { return self.id }

func (self *simple_task) Execute(vm *otto.Otto, loop *loop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		e, err := vm.Call(`new Error`, nil, self.err.Error())
		if err != nil {
			return err
		}

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}

	arguments = append(arguments, self.result)

	//arguments = append([]interface{}{self.call}, arguments...)
	/*if _, err := self.call.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}*/
	/*if _, err := vm.Call(`Function.call.call`, nil, arguments...); err != nil {
		return err
	}**/

	return nil
}

func (self *simple_task) Cancel() {

}

func SimpleTask(vm *notto.Notto, cb otto.Value, worker func(task loop.Task) (interface{}, error)) {

	task := &simple_task{call: cb}
	vm.Runloop().Add(task)

	go func() {
		defer vm.Runloop().Ready(task)
		var (
			i interface{}
			e error
			//v otto.Value
		)
		i, e = worker(task)
		task.err = e
		if i == nil {
			task.result = otto.UndefinedValue()
		} else {
			/*task.result, e = vm.ToValue(i)
			if e != nil {
				task.err = e
			}*/
		}

	}()

}
