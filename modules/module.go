package modules

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/events"
	"github.com/kildevaeld/notto/modules/fetch"
	"github.com/kildevaeld/notto/modules/fs"
	"github.com/kildevaeld/notto/modules/fsm"
	"github.com/kildevaeld/notto/modules/global"
	"github.com/kildevaeld/notto/modules/process"
	"github.com/kildevaeld/notto/modules/promise"
	"github.com/kildevaeld/notto/modules/shell"
	"github.com/kildevaeld/notto/modules/ui"
	"github.com/kildevaeld/notto/modules/util"
)

func mustError(result error) func(err error) error {
	return func(err error) error {
		if err != nil {
			return multierror.Append(result, err)
		}
		return nil
	}
}

func Define(vm *notto.Notto) error {

	var result error
	result = mustError(result)(shell.Define(vm, false))
	result = mustError(result)(process.Define(vm))
	result = mustError(result)(util.Define(vm))
	result = mustError(result)(promise.Define(vm))
	result = mustError(result)(fetch.Define(vm))
	result = mustError(result)(global.Define(vm))
	result = mustError(result)(fs.Define(vm))
	result = mustError(result)(ui.Define(vm))
	result = mustError(result)(fsm.Define(vm))
	result = mustError(result)(events.Define(vm))

	return result

}
