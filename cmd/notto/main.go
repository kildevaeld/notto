package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kildevaeld/notto/loop/looptask"
	"github.com/kildevaeld/notto/modules"

	"github.com/robertkrimen/otto"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/repl"
	"github.com/spf13/pflag"
)

var stringSlice []string
var globalShell bool

func init() {
	pflag.StringSliceVarP(&stringSlice, "env", "e", nil, "Set environment variables")
	pflag.BoolVarP(&globalShell, "shell", "s", false, "Make det shell module global")
}

func main() {

	pflag.Parse()

	vm := notto.New()

	//vm.Init()
	/*global.Define(vm)
	shell.Define(vm, globalShell)
	util.Define(vm)
	process.Define(vm)
	promise.Define(vm)
	fetch.Define(vm)
	ui.Define(vm)
	fsm.Define(vm)*/
	modules.Define(vm)
	//env := envFromStringSlice(os.Environ(), stringSlice)
	var env []string
	env = append(env, os.Environ()...)
	env = append(env, stringSlice...)

	/*process.DefineFunc(vm, func() (notto.Environ, []string) {
		return env, pflag.Args()
	})*/
	vm.ProcessAttr().Environ = env
	vm.ProcessAttr().Argv = pflag.Args()

	if len(os.Args) < 2 {
		fmt.Printf("usage: notto <path|repl>\n")
		os.Exit(1)
	}
	var v otto.Value
	var e error

	if pflag.Arg(0) == "repl" {
		e = runRepl(vm)
	} else {
		v, e = vm.Run(pflag.Arg(0), ".")
	}

	if e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		os.Exit(2)
	}

	if v.IsUndefined() || v.IsNull() {
		return
	}

}

func runRepl(vm *notto.Notto) error {
	blockingTask := looptask.NewEvalTask("")

	vm.Runloop().Add(blockingTask)
	l := vm.Runloop()
	go func() {
		if err := repl.Run(l); err != nil && err != io.EOF {
			panic(err)
		}

		l.Ready(blockingTask)
	}()

	return l.Run()
}
