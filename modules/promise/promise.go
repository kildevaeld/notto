package promise

import (
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/timers"
)

func Define(vm *notto.Notto) error {
	if v, err := vm.Get("Promise"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	if err := timers.Define(vm); err != nil {
		return err
	}

	s, err := vm.Compile("promise.js", string(MustAsset("promise.js")))
	if err != nil {
		return err
	}

	if _, err := vm.Otto.Run(s); err != nil {
		return err
	}

	return nil
}
