package main

import (
	"fmt"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/shell"
)

func main() {
	vm := notto.New()
	if err := vm.Init(); err != nil {
		fmt.Printf("ERROR %v\n", err)
	}

	shell.Define(vm, false)

	v, err := vm.Run("index.js", ".")

	fmt.Printf("%v - %v\n", v, err)

	v, _ = vm.RunScript("return 'Hello' + ;", ".")
	s, _ := v.ToString()
	fmt.Printf("Return: %s", s)
}
