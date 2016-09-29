//go:generate go-bindata -pkg docker -o docker_impl.go  builder.js docker.js
package docker

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

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

	if _, err := self.call.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}

	return nil
}

func (self *simple_task) Cancel() {

}

func SimpleTask(vm *notto.Notto, cb otto.Value, worker func() (interface{}, error)) {

	task := &simple_task{call: cb}
	vm.Runloop().Add(task)

	go func() {
		defer vm.Runloop().Ready(task)
		var (
			i interface{}
			e error
			//v otto.Value
		)
		i, e = worker()
		task.err = e
		if i == nil {
			task.result = otto.UndefinedValue()
		} else {
			task.result, e = vm.ToValue(i)
			if e != nil {
				task.err = e
			}
		}

	}()

}

func Define(vm *notto.Notto, options *DockerOptions) error {

	o, e := vm.Object("({})")
	if e != nil {
		return e
	}

	o.Set("create", func(call otto.FunctionCall) otto.Value {
		var (
			e error
			c *docker_p
			v otto.Value
		)

		if options == nil {
			options = &DockerOptions{Env: true}
		}

		if c, e = createDocker(vm, *options); e != nil {
			vm.Throw("DockerError", e)
		}
		if v, e = vm.ToValue(c); e != nil {
			vm.Throw("DockerError", e)
		}
		return v
	})

	d := make(map[string]otto.Value)
	d["__docker"] = o.Value()
	source := MustAsset("docker.js")
	loader := notto.CreateLoaderFromSourceWithPrivate(string(source), "", d)

	vm.AddModule("docker", loader)
	vm.AddModule("docker.builder", notto.CreateLoaderFromSource(string(MustAsset("builder.js")), ""))

	return nil
}

func privateDocker(vm *notto.Notto, str string, call otto.Value) {
	SimpleTask(vm, call, func() (interface{}, error) {

		cmd := exec.Command("docker")
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		conf := vm.ProcessAttr()

		cmd.Stderr = stderr
		cmd.Stdout = stdout
		cmd.Dir = conf.Cwd
		cmd.Args = strings.Split(str, " ")
		err := cmd.Run()

		if err != nil {
			return nil, fmt.Errorf("CMD: %s failed\n%s", str, stderr.Bytes())
		}

		if !cmd.ProcessState.Success() {
			return nil, errors.New(cmd.ProcessState.String())
		}

		out := strings.TrimSpace(string(stdout.Bytes()))

		if stderr.Len() > 0 {
			err = fmt.Errorf("CMD: %s failed\n%s", str, stderr.Bytes())
		}

		return out, err
	})
}
