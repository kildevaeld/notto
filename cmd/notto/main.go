package main

import (
	"fmt"
	"os"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/shell"
	"github.com/kildevaeld/notto/modules/util"
	"github.com/spf13/pflag"
)

var stringSlice []string

func init() {
	pflag.StringSliceVarP(&stringSlice, "env", "e", nil, "Set environment variables")
}

func main() {

	pflag.Parse()

	vm := notto.New()

	vm.Init()

	shell.Define(vm, false)
	util.Define(vm)

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
		fmt.Printf("usage: notto <path>\n")
		os.Exit(1)
	}

	v, e := vm.Run(pflag.Arg(0), ".")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		os.Exit(2)
	}

	if v.IsUndefined() || v.IsNull() {
		return
	}

}
