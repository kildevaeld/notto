package shell

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto, global bool) error {
	/*if v, err := vm.Get("__private_sh"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}*/

	ob, e := vm.Object("({})")
	if e != nil {
		return e
	}
	ob.Set("echo", echoCall(vm))
	ob.Set("exec", execCall(vm, "./"))
	ob.Set("cd", cdCall(vm))
	ob.Set("cat", catCall(vm))

	if global {
		vm.Set("__private_sh", ob)
		vm.AddPreScript(`var echo = __private_sh.echo;
			var exec = __private_sh.exec;
			var echo = __private_sh.echo;
			var cd = __private_sh.cd;
			var cat = __private_sh.cat;
			`)
	} else {
		notto.AddModule("shell", notto.CreateLoaderFromValue(ob.Value()))
	}

	return nil
}
