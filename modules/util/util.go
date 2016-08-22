package util

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto) error {
	if v, err := vm.Get("uril"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	ob, e := vm.Object("({})")
	if e != nil {
		return e
	}

	ob.Set("format", formatCall(vm))

	vm.AddModule("util", notto.CreateLoaderFromValue(ob.Value()))

	return nil
}
