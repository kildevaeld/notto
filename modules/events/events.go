//go:generate go-bindata -pkg events -o events_impl.go events.js
package events

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto) error {

	loader := notto.CreateLoaderFromSource(string(MustAsset("events.js")), "")

	vm.AddModule("events", loader)

	return nil
}
