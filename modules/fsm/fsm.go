//go:generate go-bindata -pkg fsm -o fsm_impl.go fsm.js
package fsm

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto) error {

	vm.AddModule("fsm", notto.CreateLoaderFromSource(string(MustAsset("fsm.js")), ""))

	return nil
}
