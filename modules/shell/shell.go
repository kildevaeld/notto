package shell

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto, global bool) error {
	if v, err := vm.Get("__shell"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	ob, e := vm.Object("({})")
	if e != nil {
		return e
	}
	ob.Set("echo", echoCall(vm))
	ob.Set("exec", execCall(vm, "./"))
	ob.Set("cd", cdCall(vm))

	if global {
		vm.Set("__sh", ob)
		vm.AddPreScript(`var echo = __sh.echo;
			var exec = __sh.exec;
			`)
	} else {
		notto.AddModule("sh", notto.CreateLoaderFromValue(ob.Value()))
	}

	return nil
}
