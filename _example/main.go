package main

import (
	"fmt"

	"github.com/kildevaeld/notto/modules/fetch"

	"github.com/kildevaeld/notto/modules/process"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/global"
	"github.com/kildevaeld/notto/modules/shell"
	"github.com/kildevaeld/notto/modules/ui"
)

func main() {
	vm := notto.New()

	global.Define(vm)
	process.Define(vm)
	fetch.Define(vm)

	shell.Define(vm, false)
	ui.Define(vm)
	v, err := vm.Run("index.js", ".")

	fmt.Printf("%v - %v\n", v, err)

	v, _ = vm.RunScript("return 'Hello' + ;", ".")
	s, _ := v.ToString()
	fmt.Printf("Return: %s", s)
}
